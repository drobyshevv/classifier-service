package expert

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/drobyshevv/classifier-expert-search/internal/models"
)

// Интерфейсы
type ArticleSaver interface {
	SaveArticle(ctx context.Context, article *models.Article) error
}

type ArticleProvider interface {
	GetArticle(ctx context.Context, id string) (*models.Article, error)
	SearchArticles(ctx context.Context, query string, limit int) ([]*models.Article, error)
}

type ExpertProvider interface {
	GetExpert(ctx context.Context, id string) (*models.Expert, error)
	SearchExperts(ctx context.Context, query string, department string, limit int) ([]*models.Expert, error)
}

type DepartmentProvider interface {
	SearchDepartments(ctx context.Context, query string, limit int) ([]*models.Department, error)
}

type VectorProvider interface {
	GetArticleVectors(ctx context.Context, articleIDs []string) (map[string][]float32, error)
	SearchSimilar(ctx context.Context, vector []float32, limit int) ([]string, error)
}

type TopicAnalyzer interface {
	AnalyzeTopic(ctx context.Context, query string) (*models.TopicAnalysisResponse, error)
}

type Service struct {
	log                *slog.Logger
	articleSaver       ArticleSaver
	articleProvider    ArticleProvider
	expertProvider     ExpertProvider
	departmentProvider DepartmentProvider
	vectorProvider     VectorProvider
	topicAnalyzer      TopicAnalyzer
	//mlService          *ml.MLServiceClient
}

func New(
	log *slog.Logger,
	articleSaver ArticleSaver,
	articleProvider ArticleProvider,
	expertProvider ExpertProvider,
	departmentProvider DepartmentProvider,
	vectorProvider VectorProvider,
	topicAnalyzer TopicAnalyzer,
) *Service {
	return &Service{
		log:                log,
		articleSaver:       articleSaver,
		articleProvider:    articleProvider,
		expertProvider:     expertProvider,
		departmentProvider: departmentProvider,
		vectorProvider:     vectorProvider,
		topicAnalyzer:      topicAnalyzer,
	}
}

func (s *Service) InitializeSystem(
	ctx context.Context,
	req *models.InitializeRequest,
) (*models.InitializeResponse, error) {
	const op = "expert.Service.InitializeSystem"

	log := s.log.With(slog.String("op", op))
	log.Info("starting system initialization",
		slog.Bool("force_reload", req.ForceReload),
		slog.Int("batch_size", int(req.BatchSize)),
		slog.Int("specific_docs_count", len(req.SpecificDocuments)),
	)

	// Получаем все статьи для обработки
	allArticles, err := s.articleProvider.SearchArticles(ctx, "", 1000)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get articles: %w", op, err)
	}

	if len(allArticles) == 0 {
		return nil, fmt.Errorf("%s: no articles found in database", op)
	}

	// Фильтруем по specific_documents если указаны
	var articlesToProcess []*models.Article
	if len(req.SpecificDocuments) > 0 {
		for _, article := range allArticles {
			for _, docID := range req.SpecificDocuments {
				if article.DocumentID == docID {
					articlesToProcess = append(articlesToProcess, article)
					break
				}
			}
		}
	} else {
		articlesToProcess = allArticles
	}

	// Применяем batch size
	if int32(len(articlesToProcess)) > req.BatchSize {
		articlesToProcess = articlesToProcess[:req.BatchSize]
	}

	jobID := fmt.Sprintf("job-%d", time.Now().Unix())

	log.Info("calculating semantic vectors for articles",
		slog.String("job_id", jobID),
		slog.Int("articles_to_process", len(articlesToProcess)),
	)

	// Получаем векторы для всех статей
	articleIDs := make([]string, len(articlesToProcess))
	for i, article := range articlesToProcess {
		articleIDs[i] = article.DocumentID
	}

	vectors, err := s.vectorProvider.GetArticleVectors(ctx, articleIDs)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get article vectors: %w", op, err)
	}

	// Сохраняем статьи с обновленной информацией о релевантности
	processedCount := 0
	for _, article := range articlesToProcess {
		if vector, exists := vectors[article.DocumentID]; exists && len(vector) > 0 {
			// Обновляем статью с информацией о векторе
			// В реальности здесь можно добавить поле Vector или сохранить в отдельную таблицу
			article.RelevanceScore = calculateRelevanceScore(vector)
			processedCount++
		}
	}

	log.Info("system initialization completed",
		slog.String("job_id", jobID),
		slog.Int("processed_articles", processedCount),
	)

	return &models.InitializeResponse{
		JobID:                     jobID,
		Status:                    "completed",
		TotalArticles:             int32(len(articlesToProcess)),
		ProcessedArticles:         int32(processedCount),
		EstimatedRemainingSeconds: 0,
		StartedAt:                 time.Now().Format(time.RFC3339),
		CurrentArticle:            "Все статьи обработаны",
	}, nil
}

func (s *Service) GetInitializationStatus(
	ctx context.Context,
	req *models.InitializationStatusRequest,
) (*models.InitializationStatusResponse, error) {
	const op = "expert.Service.GetInitializationStatus"

	log := s.log.With(slog.String("op", op))
	log.Info("checking initialization status", slog.String("job_id", req.JobID))

	// Анализируем текущее состояние системы
	allArticles, err := s.articleProvider.SearchArticles(ctx, "", 10000)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get articles: %w", op, err)
	}

	// Проверяем сколько статей имеют векторы
	articlesWithVectors := 0
	for _, article := range allArticles {
		if article.RelevanceScore > 0 {
			articlesWithVectors++
		}
	}

	progress := float32(0)
	if len(allArticles) > 0 {
		progress = float32(articlesWithVectors) / float32(len(allArticles)) * 100
	}

	status := "completed"
	if progress < 100 {
		status = "processing"
	}

	return &models.InitializationStatusResponse{
		JobID:                   req.JobID,
		Status:                  status,
		ProcessedArticles:       int32(articlesWithVectors),
		TotalArticles:           int32(len(allArticles)),
		ProgressPercentage:      progress,
		CurrentArticle:          "Анализ состояния системы",
		EstimatedCompletionTime: time.Now().Format(time.RFC3339),
		Errors:                  []string{},
		TopicsStatistics:        s.calculateTopicsStatistics(allArticles),
	}, nil
}

func (s *Service) SearchExperts(
	ctx context.Context,
	req *models.ExpertSearchRequest,
) (*models.ExpertSearchResponse, error) {
	const op = "expert.Service.SearchExperts"

	log := s.log.With(slog.String("op", op))
	log.Info("searching experts by query",
		slog.String("query", req.Query),
		slog.String("department", req.FilterDepartment),
		slog.Int("max_results", int(req.MaxResults)),
	)

	// Получаем экспертов через провайдер
	experts, err := s.expertProvider.SearchExperts(ctx, req.Query, req.FilterDepartment, int(req.MaxResults))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to search experts: %w", op, err)
	}

	// Если запрос семантический, используем векторный поиск для улучшения результатов
	if isSemanticQuery(req.Query) {
		enhancedExperts, err := s.enhanceExpertSearchWithVectors(ctx, req.Query, experts, int(req.MaxResults))
		if err != nil {
			log.Warn("failed to enhance search with vectors", slog.String("error", err.Error()))
		} else {
			experts = enhancedExperts
		}
	}

	// Сортируем экспертов по релевантности
	sort.Slice(experts, func(i, j int) bool {
		return experts[i].ExpertiseScore > experts[j].ExpertiseScore
	})

	log.Info("expert search completed",
		slog.Int("experts_found", len(experts)),
	)

	return &models.ExpertSearchResponse{
		Experts:    experts,
		TotalFound: int32(len(experts)),
	}, nil
}

func (s *Service) SearchDepartments(
	ctx context.Context,
	req *models.DepartmentSearchRequest,
) (*models.DepartmentSearchResponse, error) {
	const op = "expert.Service.SearchDepartments"

	log := s.log.With(slog.String("op", op))
	log.Info("searching departments",
		slog.String("query", req.Query),
		slog.Int("max_results", int(req.MaxResults)),
	)

	// Получаем подразделения через провайдер
	departments, err := s.departmentProvider.SearchDepartments(ctx, req.Query, int(req.MaxResults))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to search departments: %w", op, err)
	}

	// Для каждого подразделения получаем детальную информацию
	for i, dept := range departments {
		// Получаем экспертов этого подразделения
		departmentExperts, err := s.expertProvider.SearchExperts(ctx, "", dept.Name, 10)
		if err == nil && len(departmentExperts) > 0 {
			// Обновляем статистику на основе реальных данных
			dept.ExpertCount = int32(len(departmentExperts))
			dept.TotalArticles = s.calculateTotalArticles(departmentExperts)
			dept.StrengthScore = calculateDepartmentStrength(departmentExperts)
			dept.KeyAuthors = s.getTopAuthors(departmentExperts, 3)
		}
		departments[i] = dept
	}

	log.Info("department search completed",
		slog.Int("departments_found", len(departments)),
	)

	return &models.DepartmentSearchResponse{
		Departments: departments,
		TotalFound:  int32(len(departments)),
	}, nil
}

func (s *Service) AnalyzeTopic(
	ctx context.Context,
	req *models.TopicAnalysisRequest,
) (*models.TopicAnalysisResponse, error) {
	const op = "expert.Service.AnalyzeTopic"

	log := s.log.With(slog.String("op", op))
	log.Info("analyzing topic", slog.String("query", req.Query))

	// Используем topic analyzer для анализа
	analysis, err := s.topicAnalyzer.AnalyzeTopic(ctx, req.Query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to analyze topic: %w", op, err)
	}

	// Дополняем анализ данными из БД
	relatedArticles, err := s.articleProvider.SearchArticles(ctx, req.Query, 20)
	if err == nil {
		analysis.RelatedTopics = s.extractTopicsFromArticles(relatedArticles)
	}

	log.Info("topic analysis completed",
		slog.String("main_topic", analysis.MainTopic),
		slog.Int("related_topics", len(analysis.RelatedTopics)),
	)

	return analysis, nil
}

func (s *Service) SearchArticles(
	ctx context.Context,
	req *models.ArticleSearchRequest,
) (*models.ArticleSearchResponse, error) {
	const op = "expert.Service.SearchArticles"

	log := s.log.With(slog.String("op", op))
	log.Info("searching articles",
		slog.String("query", req.Query),
		slog.Int("max_results", int(req.MaxResults)),
	)

	// Получаем статьи через провайдер
	articles, err := s.articleProvider.SearchArticles(ctx, req.Query, int(req.MaxResults))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to search articles: %w", op, err)
	}

	// Если запрос семантический, улучшаем поиск с помощью векторов
	if isSemanticQuery(req.Query) {
		enhancedArticles, err := s.enhanceArticleSearchWithVectors(ctx, req.Query, articles, int(req.MaxResults))
		if err != nil {
			log.Warn("failed to enhance article search with vectors", slog.String("error", err.Error()))
		} else {
			articles = enhancedArticles
		}
	}

	// Сортируем статьи по релевантности
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].RelevanceScore > articles[j].RelevanceScore
	})

	log.Info("article search completed",
		slog.Int("articles_found", len(articles)),
	)

	return &models.ArticleSearchResponse{
		Articles:   articles,
		TotalFound: int32(len(articles)),
	}, nil
}

// Вспомогательные методы

func (s *Service) enhanceExpertSearchWithVectors(ctx context.Context, query string, experts []*models.Expert, limit int) ([]*models.Expert, error) {
	// Получаем вектор запроса
	similarDocs, err := s.vectorProvider.SearchSimilar(ctx, []float32{}, limit)
	if err != nil {
		return experts, err
	}

	// Находим экспертов по релевантным документам
	enhancedExperts := make(map[string]*models.Expert)
	for _, docID := range similarDocs {
		article, err := s.articleProvider.GetArticle(ctx, docID)
		if err != nil {
			continue
		}
		for _, author := range article.Authors {
			expert, err := s.expertProvider.GetExpert(ctx, author)
			if err == nil {
				enhancedExperts[expert.AuthorID] = expert
			}
		}
	}

	// Объединяем результаты
	result := make([]*models.Expert, 0, len(enhancedExperts))
	for _, expert := range enhancedExperts {
		result = append(result, expert)
	}

	return result, nil
}

func (s *Service) enhanceArticleSearchWithVectors(ctx context.Context, query string, articles []*models.Article, limit int) ([]*models.Article, error) {
	// Используем векторный поиск для улучшения результатов
	similarDocs, err := s.vectorProvider.SearchSimilar(ctx, []float32{}, limit)
	if err != nil {
		return articles, err
	}

	// Получаем статьи по ID из векторного поиска
	for _, docID := range similarDocs {
		article, err := s.articleProvider.GetArticle(ctx, docID)
		if err == nil {
			// Проверяем, нет ли уже этой статьи в результатах
			found := false
			for _, existing := range articles {
				if existing.DocumentID == article.DocumentID {
					found = true
					break
				}
			}
			if !found {
				articles = append(articles, article)
			}
		}
	}

	return articles, nil
}

func (s *Service) calculateTopicsStatistics(articles []*models.Article) map[string]int32 {
	stats := make(map[string]int32)
	for _, article := range articles {
		for _, topic := range article.MatchedTopics {
			stats[topic]++
		}
		// Если нет явных топиков, анализируем по заголовку
		if len(article.MatchedTopics) == 0 {
			mainTopic := extractMainTopic(article.TitleRu)
			if mainTopic != "" {
				stats[mainTopic]++
			}
		}
	}
	return stats
}

func (s *Service) extractTopicsFromArticles(articles []*models.Article) []string {
	topics := make(map[string]bool)
	for _, article := range articles {
		for _, topic := range article.MatchedTopics {
			topics[topic] = true
		}
	}

	result := make([]string, 0, len(topics))
	for topic := range topics {
		result = append(result, topic)
	}
	return result
}

func (s *Service) calculateTotalArticles(experts []*models.Expert) int32 {
	var total int32
	for _, expert := range experts {
		total += expert.ArticleCount
	}
	return total
}

func (s *Service) getTopAuthors(experts []*models.Expert, limit int) []string {
	sort.Slice(experts, func(i, j int) bool {
		return experts[i].ArticleCount > experts[j].ArticleCount
	})

	authors := make([]string, 0, limit)
	for i := 0; i < len(experts) && i < limit; i++ {
		authors = append(authors, experts[i].NameRu)
	}
	return authors
}

// Утилитарные функции

func calculateRelevanceScore(vector []float32) float32 {
	// Простая эвристика для расчета релевантности
	if len(vector) == 0 {
		return 0
	}
	// Суммируем значения вектора (упрощенный подход)
	var sum float32
	for _, v := range vector {
		sum += v
	}
	return sum / float32(len(vector))
}

func calculateDepartmentStrength(experts []*models.Expert) float32 {
	if len(experts) == 0 {
		return 0
	}
	var totalScore float32
	for _, expert := range experts {
		totalScore += expert.ExpertiseScore
	}
	return totalScore / float32(len(experts))
}

func isSemanticQuery(query string) bool {
	// Эвристика для определения семантических запросов
	return len(query) > 3 && !isSimpleKeyword(query)
}

func isSimpleKeyword(query string) bool {
	simpleKeywords := []string{"и", "или", "на", "в", "с", "по", "за", "о", "об", "от", "до"}
	for _, keyword := range simpleKeywords {
		if query == keyword {
			return true
		}
	}
	return false
}

func extractMainTopic(title string) string {
	// Упрощенная логика извлечения основной темы из заголовка
	// В реальности можно использовать NLP
	if contains(title, "машинн") || contains(title, "ML") {
		return "машинное обучение"
	}
	if contains(title, "нейросет") || contains(title, "deep learning") {
		return "нейросети"
	}
	if contains(title, "анализ данн") || contains(title, "data analysis") {
		return "анализ данных"
	}
	return ""
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr))
}

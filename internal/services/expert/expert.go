package expert

import (
	"context"
	"encoding/binary"
	"fmt"
	"log/slog"
	"math"
	"sort"
	"time"

	"github.com/drobyshevv/classifier-expert-search/internal/models"
	"github.com/drobyshevv/classifier-expert-search/internal/services/ml"
	agentv1 "github.com/drobyshevv/proto-ai-agent/gen/go/proto/ai_agent"
)

// Интерфейсы
type ArticleSaver interface {
	SaveArticle(ctx context.Context, article *models.Article) error
}

type ArticleProvider interface {
	GetArticle(ctx context.Context, id string) (*models.Article, error)
	SearchArticles(ctx context.Context, query string, limit int) ([]*models.Article, error)
	SearchArticlesPaged(ctx context.Context, limit, offset int) ([]models.Article, error)
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
	SaveArticleEmbedding(ctx context.Context, documentID string, titleEmbedding, abstractEmbedding []float32) error
	HasArticleEmbedding(ctx context.Context, documentID string) (bool, error)

	SaveArticleTopics(ctx context.Context, documentID string, topics []models.ArticleTopic) error
	RecalculateAggregatesForDocument(ctx context.Context, documentID string) error
	ClearEmbeddingsAndTopics(ctx context.Context) error
	ListArticlesForSemantic(ctx context.Context, limit int) ([]models.SemanticArticle, error)
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
	mlService          *ml.MLServiceClient
}

func New(
	log *slog.Logger,
	articleSaver ArticleSaver,
	articleProvider ArticleProvider,
	expertProvider ExpertProvider,
	departmentProvider DepartmentProvider,
	vectorProvider VectorProvider,
	topicAnalyzer TopicAnalyzer,
	mlService *ml.MLServiceClient,
) *Service {
	return &Service{
		log:                log,
		articleSaver:       articleSaver,
		articleProvider:    articleProvider,
		expertProvider:     expertProvider,
		departmentProvider: departmentProvider,
		vectorProvider:     vectorProvider,
		topicAnalyzer:      topicAnalyzer,
		mlService:          mlService,
	}
}

/*
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
*/
/* v2
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

	log.Info("generating semantic vectors for articles",
		slog.String("job_id", jobID),
		slog.Int("articles_to_process", len(articlesToProcess)),
	)

	// ГЕНЕРИРУЕМ ВЕКТОРЫ ЧЕРЕЗ ML СЕРВИС
	processedCount := 0
	for _, article := range articlesToProcess {
		// Вызываем AI Agent для генерации эмбеддингов
		analysisResult, err := s.mlService.AnalyzeArticleTopics(ctx, article.DocumentID, article.TitleRu, article.AbstractRu)
		if err != nil {
			log.Warn("failed to generate embeddings for article",
				slog.String("document_id", article.DocumentID),
				slog.String("error", err.Error()))
			continue
		}

		log.Info("DEBUG: Received analysis result",
			slog.String("document_id", article.DocumentID),
			slog.Int("title_embedding_len", len(analysisResult.TitleEmbedding)),
			slog.Int("abstract_embedding_len", len(analysisResult.AbstractEmbedding)),
			slog.String("title_preview", string(analysisResult.TitleEmbedding[:min(50, len(analysisResult.TitleEmbedding))])),
			slog.String("abstract_preview", string(analysisResult.AbstractEmbedding[:min(50, len(analysisResult.AbstractEmbedding))])),
		)

		// Декодируем base64 эмбеддинги в float32
		titleVector, err := bytesToFloat32(analysisResult.TitleEmbedding)
		if err != nil {
			log.Warn("failed to decode title embedding",
				slog.String("document_id", article.DocumentID),
				slog.String("error", err.Error()))
			continue
		}

		abstractVector, err := bytesToFloat32(analysisResult.AbstractEmbedding)
		if err != nil {
			log.Warn("failed to decode abstract embedding",
				slog.String("document_id", article.DocumentID),
				slog.String("error", err.Error()))
			continue
		}

		// СОХРАНЯЕМ ВЕКТОРЫ В БАЗУ
		err = s.vectorProvider.SaveArticleEmbedding(ctx, article.DocumentID, titleVector, abstractVector)
		if err != nil {
			log.Warn("failed to save article embedding",
				slog.String("document_id", article.DocumentID),
				slog.String("error", err.Error()))
			continue
		}

		// Обновляем статью с информацией о релевантности
		article.RelevanceScore = calculateRelevanceScore(titleVector)

		// Сохраняем тематики
		var matchedTopics []string
		for _, topic := range analysisResult.Topics {
			matchedTopics = append(matchedTopics, topic.TopicName)
		}
		article.MatchedTopics = matchedTopics

		processedCount++
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

*/

//TODO: обработать documents
//TODO: безопасные батчи/распараллелевание
//FIXME: Нулевые векторы - ПРОБЛЕМА!  Добавить проверку в коде
// добавить проверку на нулевые векторы чтобы не засорять базу бесполезными данными!

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

	// batch size default
	batchSize := 50
	if req.BatchSize > 0 && req.BatchSize < 50 {
		batchSize = int(req.BatchSize)
	}

	// force reload: очистка таблиц
	if req.ForceReload {
		if err := s.vectorProvider.ClearEmbeddingsAndTopics(ctx); err != nil {
			log.Warn("failed to clear embeddings/topics", slog.String("error", err.Error()))
		} else {
			log.Info("cleared embeddings and topics due to force_reload")
		}
	}

	offset := 0
	totalProcessed := 0
	hasMore := true

	jobID := fmt.Sprintf("job-%d", time.Now().Unix())

	log.Info("starting batch processing",
		slog.String("job_id", jobID),
		slog.Int("batch_size", batchSize),
	)

	for hasMore {
		log.Info("processing batch",
			slog.Int("batch_number", offset/batchSize+1),
			slog.Int("offset", offset),
		)

		// get batch (use storage method which returns limited number)
		batchArticles, err := s.articleProvider.SearchArticlesPaged(ctx, batchSize, offset)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get articles batch: %w", op, err)
		}

		if len(batchArticles) == 0 {
			hasMore = false
			log.Info("no more articles to process")
			break
		}

		batchProcessed := 0
		for i, article := range batchArticles {
			// progress log
			if i%10 == 0 {
				log.Info("batch progress",
					slog.Int("batch_processed", i),
					slog.Int("batch_total", len(batchArticles)),
					slog.Int("total_processed", totalProcessed),
				)
			}

			// skip existing if asked
			if req.SkipExisting {
				hasVector, herr := s.vectorProvider.HasArticleEmbedding(ctx, article.DocumentID)
				if herr == nil && hasVector {
					log.Debug("skipping article with existing embedding",
						slog.String("document_id", article.DocumentID))
					continue
				}
			}

			// call ML
			analysisResult, err := s.mlService.AnalyzeArticleTopics(ctx, article.DocumentID, article.TitleRu, article.AbstractRu)
			if err != nil {
				log.Warn("failed to generate embeddings for article",
					slog.String("document_id", article.DocumentID),
					slog.String("error", err.Error()))
				continue
			}

			// decode bytes -> []float32 (service bytesToFloat32 helper)
			titleVector, err := bytesToFloat32(analysisResult.TitleEmbedding)
			if err != nil {
				log.Warn("failed to decode title embedding",
					slog.String("document_id", article.DocumentID),
					slog.String("error", err.Error()))
				continue
			}
			abstractVector, err := bytesToFloat32(analysisResult.AbstractEmbedding)
			if err != nil {
				log.Warn("failed to decode abstract embedding",
					slog.String("document_id", article.DocumentID),
					slog.String("error", err.Error()))
				continue
			}

			// Проверка: хотя бы один вектор не нулевой
			titleZero := isZeroVector(titleVector)
			abstractZero := isZeroVector(abstractVector)

			if titleZero && abstractZero {
				log.Warn("skipping article because both embeddings are zero",
					slog.String("document_id", article.DocumentID),
				)
				continue
			}

			// Если title пустой — используем abstract вместо него
			if titleZero && !abstractZero {
				titleVector = abstractVector
				log.Info("title embedding was zero, replaced with abstract embedding",
					slog.String("document_id", article.DocumentID))
			}

			// Если abstract пустой — используем title вместо него
			if abstractZero && !titleZero {
				abstractVector = titleVector
				log.Info("abstract embedding was zero, replaced with title embedding",
					slog.String("document_id", article.DocumentID))
			}

			// save embeddings
			if err := s.vectorProvider.SaveArticleEmbedding(ctx, article.DocumentID, titleVector, abstractVector); err != nil {
				log.Warn("failed to save article embedding",
					slog.String("document_id", article.DocumentID),
					slog.String("error", err.Error()))
				continue
			}

			// save topics (if any)
			if len(analysisResult.Topics) > 0 {
				var topics []models.ArticleTopic
				for _, t := range analysisResult.Topics {
					topics = append(topics, models.ArticleTopic{
						TopicName:  t.TopicName,
						Confidence: t.Confidence,
						TopicType:  t.TopicType,
					})
				}
				if err := s.vectorProvider.SaveArticleTopics(ctx, article.DocumentID, topics); err != nil {
					log.Warn("failed to save article topics",
						slog.String("document_id", article.DocumentID),
						slog.String("error", err.Error()))
				} else {
					// recalc aggregates (topic_experts, department_topics)
					if err := s.vectorProvider.RecalculateAggregatesForDocument(ctx, article.DocumentID); err != nil {
						log.Warn("failed to recalc aggregates",
							slog.String("document_id", article.DocumentID),
							slog.String("error", err.Error()))
					}
				}
			}

			batchProcessed++
			totalProcessed++
		}

		log.Info("batch completed",
			slog.Int("batch_number", offset/batchSize+1),
			slog.Int("batch_processed", batchProcessed),
			slog.Int("total_processed", totalProcessed),
		)

		// next slice - currently using SearchArticles with limit (so increment offset by batchSize)
		offset += batchSize

		// small sleep
		time.Sleep(2 * time.Second)
	}

	totalArticles, err := s.getTotalArticlesCount(ctx)
	if err != nil {
		log.Warn("failed to get total articles count", slog.String("error", err.Error()))
		totalArticles = totalProcessed
	}

	log.Info("system initialization completed",
		slog.String("job_id", jobID),
		slog.Int("total_processed", totalProcessed),
		slog.Int("total_articles", totalArticles),
	)

	return &models.InitializeResponse{
		JobID:                     jobID,
		Status:                    "completed",
		TotalArticles:             int32(totalArticles),
		ProcessedArticles:         int32(totalProcessed),
		EstimatedRemainingSeconds: 0,
		StartedAt:                 time.Now().Format(time.RFC3339),
		CurrentArticle:            fmt.Sprintf("Обработано %d из %d статей", totalProcessed, totalArticles),
	}, nil
}

// Вспомогательные методы

func isZeroVector(v []float32) bool {
	if len(v) == 0 {
		return true
	}
	for _, x := range v {
		if x != 0 {
			return false
		}
	}
	return true
}

func (s *Service) getArticlesBatch(ctx context.Context, offset, limit int) ([]*models.Article, error) {
	// Используем существующий метод SearchArticles с большим лимитом и фильтрацией по offset
	allArticles, err := s.articleProvider.SearchArticles(ctx, "", 100000)
	if err != nil {
		return nil, err
	}

	// Ручная пагинация (в идеале добавить пагинацию в storage)
	if offset >= len(allArticles) {
		return []*models.Article{}, nil
	}

	end := offset + limit
	if end > len(allArticles) {
		end = len(allArticles)
	}

	return allArticles[offset:end], nil
}

func (s *Service) getTotalArticlesCount(ctx context.Context) (int, error) {
	articles, err := s.articleProvider.SearchArticles(ctx, "", 100000)
	if err != nil {
		return 0, err
	}
	return len(articles), nil
}

func bytesToFloat32(data []byte) ([]float32, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	if len(data)%4 != 0 {
		return nil, fmt.Errorf("data length %d not divisible by 4", len(data))
	}

	vector := make([]float32, len(data)/4)
	for i := 0; i < len(vector); i++ {
		bits := binary.LittleEndian.Uint32(data[i*4 : (i+1)*4])
		vector[i] = math.Float32frombits(bits)
	}

	return vector, nil
}

func (s *Service) SemanticSearchArticles(
	ctx context.Context,
	req *models.SemanticArticleSearchRequest,
) (*models.SemanticArticleSearchResponse, error) {

	log := s.log.With(slog.String("op", "SemanticSearchArticles"))
	log.Info("semantic search articles", slog.String("query", req.Query))

	desired := req.MaxResults
	if desired <= 0 {
		desired = 10
	}

	// -----------------------------
	// 1. Анализ запроса
	// -----------------------------
	qResp, err := s.mlService.AnalyzeUserQuery(ctx, req.Query, "article_search")
	if err != nil {
		return nil, fmt.Errorf("failed to analyze query: %w", err)
	}

	if len(qResp.QueryVector) == 0 {
		return nil, fmt.Errorf("empty query vector from ML")
	}

	// -----------------------------
	// 2. Получаем статьи
	// -----------------------------
	limit := desired + req.Offset
	items, err := s.vectorProvider.ListArticlesForSemantic(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed list articles: %w", err)
	}

	if len(items) == 0 {
		return &models.SemanticArticleSearchResponse{}, nil
	}

	// -----------------------------
	// 3. Готовим proto ArticleForSearch
	// -----------------------------
	protoArticles := make([]*agentv1.ArticleForSearch, 0, len(items))

	for _, a := range items {
		protoArticles = append(protoArticles, &agentv1.ArticleForSearch{
			DocumentId:        a.DocumentID,
			TitleRu:           a.TitleRu,
			AbstractRu:        a.AbstractRu,
			TitleEmbedding:    a.TitleEmbedding,    // []byte OK
			AbstractEmbedding: a.AbstractEmbedding, // []byte OK
		})
		log.Info("emb_len", slog.Int("title", len(a.TitleEmbedding)), slog.Int("abstract", len(a.AbstractEmbedding)))

	}

	// -----------------------------
	// 4. Proto запрос в AI Agent
	// -----------------------------
	protoReq := &agentv1.SemanticSearchRequest{
		QueryVector: qResp.QueryVector, // []byte
		Articles:    protoArticles,
		MaxResults:  int32(desired),
	}

	protoResp, err := s.mlService.SemanticArticleSearch(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("python semantic search failed: %w", err)
	}

	// -----------------------------
	// 5. Собираем результат
	// -----------------------------
	results := make([]models.SemanticArticleResult, 0, len(protoResp.Results))

	fastMeta := map[string]models.SemanticArticle{}
	for _, a := range items {
		fastMeta[a.DocumentID] = a
	}

	for _, r := range protoResp.Results {
		m := fastMeta[r.DocumentId]

		art, _ := s.articleProvider.GetArticle(ctx, r.DocumentId)

		results = append(results, models.SemanticArticleResult{
			DocumentID:     r.DocumentId,
			TitleRu:        m.TitleRu,
			AbstractRu:     m.AbstractRu,
			Year:           art.Year,
			Authors:        art.Authors,
			RelevanceScore: r.RelevanceScore,
			MatchedTopics:  r.MatchedConcepts,
		})
	}

	return &models.SemanticArticleSearchResponse{
		Results:    results,
		TotalFound: int(protoResp.TotalFound),
	}, nil
}

func float32ToBytes(v []float32) []byte {
	buf := make([]byte, len(v)*4)
	for i, f := range v {
		binary.LittleEndian.PutUint32(buf[i*4:(i+1)*4], math.Float32bits(f))
	}
	return buf
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

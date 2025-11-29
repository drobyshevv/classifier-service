package expertgrpc

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/drobyshevv/classifier-expert-search/internal/models"
	servicev1 "github.com/drobyshevv/proto-classifier-expert-search/gen/go/proto/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Expert interface {
	InitializeSystem(ctx context.Context, req *models.InitializeRequest) (*models.InitializeResponse, error)

	GetInitializationStatus(ctx context.Context, req *models.InitializationStatusRequest) (*models.InitializationStatusResponse, error)

	SearchExperts(ctx context.Context, req *models.ExpertSearchRequest) (*models.ExpertSearchResponse, error)

	SearchDepartments(ctx context.Context, req *models.DepartmentSearchRequest) (*models.DepartmentSearchResponse, error)

	AnalyzeTopic(ctx context.Context, req *models.TopicAnalysisRequest) (*models.TopicAnalysisResponse, error)

	SearchArticles(ctx context.Context, req *models.ArticleSearchRequest) (*models.ArticleSearchResponse, error)

	SemanticSearchArticles(ctx context.Context, req *models.SemanticArticleSearchRequest) (*models.SemanticArticleSearchResponse, error)
}

type serverAPI struct {
	servicev1.UnimplementedExpertSearchServiceServer
	expert Expert
}

func Register(gRPC *grpc.Server, expert Expert) {
	servicev1.RegisterExpertSearchServiceServer(gRPC, &serverAPI{expert: expert})
}

const (
	emptyValue = 0
)

func (s *serverAPI) InitializeSystem(ctx context.Context, req *servicev1.InitializeRequest) (*servicev1.InitializeResponse, error) {

	if err := validateInitializeRequest(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	modelReq := &models.InitializeRequest{
		ForceReload:       req.GetForceReload(),
		SkipExisting:      req.GetSkipExisting(),
		BatchSize:         req.GetBatchSize(),
		SpecificDocuments: req.GetSpecificDocuments(),
	}

	modelResp, err := s.expert.InitializeSystem(ctx, modelReq)
	if err != nil {
		return nil, err
	}

	return &servicev1.InitializeResponse{
		JobId:                     modelResp.JobID,
		Status:                    modelResp.Status,
		TotalArticles:             modelResp.TotalArticles,
		ProcessedArticles:         modelResp.ProcessedArticles,
		EstimatedRemainingSeconds: modelResp.EstimatedRemainingSeconds,
		StartedAt:                 modelResp.StartedAt,
		CurrentArticle:            modelResp.CurrentArticle,
	}, nil
}

func (s *serverAPI) SemanticSearchArticles(ctx context.Context, req *servicev1.SemanticArticleSearchRequest) (*servicev1.SemanticArticleSearchResponse, error) {
	// validate
	q := req.GetQuery()
	if q == "" {
		return nil, status.Error(codes.InvalidArgument, "query cannot be empty")
	}
	if utf8.RuneCountInString(q) > 500 {
		return nil, status.Error(codes.InvalidArgument, "query too long (max 500 chars)")
	}
	max := int(req.GetMaxResults())
	if max <= 0 {
		max = 10
	}
	if max > 200 {
		return nil, status.Error(codes.InvalidArgument, "max_results cannot exceed 200")
	}

	modelReq := &models.SemanticArticleSearchRequest{
		Query:      q,
		MaxResults: max,
		Offset:     int(req.GetOffset()),
	}

	modelResp, err := s.expert.SemanticSearchArticles(ctx, modelReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "semantic search failed: %v", err)
	}

	// map response
	results := make([]*servicev1.Article, 0, len(modelResp.Results))
	for _, a := range modelResp.Results {
		results = append(results, &servicev1.Article{
			DocumentId:     a.DocumentID,
			TitleRu:        a.TitleRu,
			AbstractRu:     a.AbstractRu,
			Year:           a.Year,
			Authors:        a.Authors,
			RelevanceScore: a.RelevanceScore,
			MatchedTopics:  a.MatchedTopics,
		})
	}

	return &servicev1.SemanticArticleSearchResponse{
		Results:    results,
		TotalFound: int32(modelResp.TotalFound),
	}, nil
}

func (s *serverAPI) GetInitializationStatus(
	ctx context.Context,
	req *servicev1.InitializationStatusRequest,
) (*servicev1.InitializationStatusResponse, error) {

	modelReq := &models.InitializationStatusRequest{
		JobID: req.GetJobId(),
	}

	modelResp, err := s.expert.GetInitializationStatus(ctx, modelReq)
	if err != nil {
		return nil, err
	}

	// Map for topics statistics
	topicsStats := make(map[string]int32)
	for k, v := range modelResp.TopicsStatistics {
		topicsStats[k] = v
	}

	return &servicev1.InitializationStatusResponse{
		JobId:                   modelResp.JobID,
		Status:                  modelResp.Status,
		ProcessedArticles:       modelResp.ProcessedArticles,
		TotalArticles:           modelResp.TotalArticles,
		ProgressPercentage:      modelResp.ProgressPercentage,
		CurrentArticle:          modelResp.CurrentArticle,
		EstimatedCompletionTime: modelResp.EstimatedCompletionTime,
		Errors:                  modelResp.Errors,
		TopicsStatistics:        topicsStats,
	}, nil
}

func (s *serverAPI) SearchExperts(
	ctx context.Context,
	req *servicev1.ExpertSearchRequest,
) (*servicev1.ExpertSearchResponse, error) {

	if err := validateExpertSearchRequest(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid search request: %v", err)
	}

	modelReq := &models.ExpertSearchRequest{
		Query:            req.GetQuery(),
		FilterDepartment: req.GetFilterDepartment(),
		MaxResults:       req.GetMaxResults(),
	}

	modelResp, err := s.expert.SearchExperts(ctx, modelReq)
	if err != nil {
		return nil, err
	}

	experts := make([]*servicev1.Expert, len(modelResp.Experts))
	for i, expert := range modelResp.Experts {
		experts[i] = &servicev1.Expert{
			AuthorId:         expert.AuthorID,
			NameRu:           expert.NameRu,
			Department:       expert.Department,
			ExpertiseScore:   expert.ExpertiseScore,
			ArticleCount:     expert.ArticleCount,
			TotalCitations:   expert.TotalCitations,
			LastActivityYear: expert.LastActivityYear,
		}
	}

	return &servicev1.ExpertSearchResponse{
		Experts:    experts,
		TotalFound: modelResp.TotalFound,
	}, nil
}

func (s *serverAPI) SearchDepartments(
	ctx context.Context,
	req *servicev1.DepartmentSearchRequest,
) (*servicev1.DepartmentSearchResponse, error) {

	if err := validateDepartmentSearchRequest(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid search request: %v", err)
	}

	modelReq := &models.DepartmentSearchRequest{
		Query:      req.GetQuery(),
		MaxResults: req.GetMaxResults(),
	}

	modelResp, err := s.expert.SearchDepartments(ctx, modelReq)
	if err != nil {
		return nil, err
	}

	departments := make([]*servicev1.Department, len(modelResp.Departments))
	for i, dept := range modelResp.Departments {
		departments[i] = &servicev1.Department{
			OrganizationId: dept.OrganizationID,
			Name:           dept.Name,
			StrengthScore:  dept.StrengthScore,
			ExpertCount:    dept.ExpertCount,
			TotalArticles:  dept.TotalArticles,
			KeyAuthors:     dept.KeyAuthors,
		}
	}

	return &servicev1.DepartmentSearchResponse{
		Departments: departments,
		TotalFound:  modelResp.TotalFound,
	}, nil
}

func (s *serverAPI) AnalyzeTopic(
	ctx context.Context,
	req *servicev1.TopicAnalysisRequest,
) (*servicev1.TopicAnalysisResponse, error) {

	if err := validateTopicAnalysisRequest(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid search request: %v", err)
	}

	modelReq := &models.TopicAnalysisRequest{
		Query: req.GetQuery(),
	}

	modelResp, err := s.expert.AnalyzeTopic(ctx, modelReq)
	if err != nil {
		return nil, err
	}

	return &servicev1.TopicAnalysisResponse{
		MainTopic:        modelResp.MainTopic,
		RelatedTopics:    modelResp.RelatedTopics,
		TopicDescription: modelResp.TopicDescription,
	}, nil
}

func (s *serverAPI) SearchArticles(
	ctx context.Context,
	req *servicev1.ArticleSearchRequest,
) (*servicev1.ArticleSearchResponse, error) {

	if err := validateArticleSearchRequest(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid search request: %v", err)
	}

	modelReq := &models.ArticleSearchRequest{
		Query:      req.GetQuery(),
		MaxResults: req.GetMaxResults(),
	}

	modelResp, err := s.expert.SearchArticles(ctx, modelReq)
	if err != nil {
		return nil, err
	}

	// Конвертируем список статей
	articles := make([]*servicev1.Article, len(modelResp.Articles))
	for i, article := range modelResp.Articles {
		articles[i] = &servicev1.Article{
			DocumentId:     article.DocumentID,
			TitleRu:        article.TitleRu,
			AbstractRu:     article.AbstractRu,
			Year:           article.Year,
			Authors:        article.Authors,
			RelevanceScore: article.RelevanceScore,
			MatchedTopics:  article.MatchedTopics,
		}
	}

	return &servicev1.ArticleSearchResponse{
		Articles:   articles,
		TotalFound: modelResp.TotalFound,
	}, nil
}

func validateInitializeRequest(req *servicev1.InitializeRequest) error {

	if req.GetBatchSize() <= 0 {
		return fmt.Errorf("batch_size must be positive")
	}

	if req.GetBatchSize() > 1000 {
		return fmt.Errorf("batch_size cannot exceed 1000")
	}

	if len(req.GetSpecificDocuments()) > 100 {
		return fmt.Errorf("too many specific documents (max: 100)")
	}

	for _, docID := range req.GetSpecificDocuments() {
		if strings.TrimSpace(docID) == "" {
			return fmt.Errorf("document ID cannot be empty")
		}
	}

	return nil
}

func validateExpertSearchRequest(req *servicev1.ExpertSearchRequest) error {
	query := strings.TrimSpace(req.GetQuery())
	if query == "" {
		return status.Error(codes.InvalidArgument, "query cannot be empty")
	}
	if utf8.RuneCountInString(query) > 500 {
		return status.Error(codes.InvalidArgument, "query too long (max: 500 characters)")
	}
	if req.GetMaxResults() < 1 {
		return status.Error(codes.InvalidArgument, "max_results must be at least 1")
	}
	if req.GetMaxResults() > 100 {
		return status.Error(codes.InvalidArgument, "max_results cannot exceed 100")
	}

	if req.GetFilterDepartment() != "" && len(req.GetFilterDepartment()) > 200 {
		return status.Error(codes.InvalidArgument, "filter_department too long (max: 200 characters)")
	}
	return nil
}

func validateDepartmentSearchRequest(req *servicev1.DepartmentSearchRequest) error {
	query := strings.TrimSpace(req.GetQuery())
	if query == "" {
		return status.Error(codes.InvalidArgument, "query cannot be empty")
	}
	if utf8.RuneCountInString(query) > 300 {
		return status.Error(codes.InvalidArgument, "query too long (max: 300 characters)")
	}
	if req.GetMaxResults() < 1 {
		return status.Error(codes.InvalidArgument, "max_results must be at least 1")
	}
	if req.GetMaxResults() > 50 {
		return status.Error(codes.InvalidArgument, "max_results cannot exceed 50")
	}
	return nil
}

func validateTopicAnalysisRequest(req *servicev1.TopicAnalysisRequest) error {
	query := strings.TrimSpace(req.GetQuery())
	if query == "" {
		return status.Error(codes.InvalidArgument, "query cannot be empty")
	}
	if utf8.RuneCountInString(query) < 3 {
		return status.Error(codes.InvalidArgument, "query too short (min: 3 characters)")
	}
	if utf8.RuneCountInString(query) > 1000 {
		return status.Error(codes.InvalidArgument, "query too long (max: 1000 characters)")
	}
	return nil
}

func validateArticleSearchRequest(req *servicev1.ArticleSearchRequest) error {
	query := strings.TrimSpace(req.GetQuery())
	if query == "" {
		return status.Error(codes.InvalidArgument, "query cannot be empty")
	}
	if utf8.RuneCountInString(query) > 500 {
		return status.Error(codes.InvalidArgument, "query too long (max: 500 characters)")
	}
	if req.GetMaxResults() < 1 {
		return status.Error(codes.InvalidArgument, "max_results must be at least 1")
	}
	if req.GetMaxResults() > 200 {
		return status.Error(codes.InvalidArgument, "max_results cannot exceed 200")
	}
	return nil
}

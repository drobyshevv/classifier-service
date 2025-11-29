package models

// REQUEST STRUCTURES
type InitializeRequest struct {
	ForceReload       bool     `json:"force_reload"`
	SkipExisting      bool     `json:"skip_existing"`
	BatchSize         int32    `json:"batch_size"`
	SpecificDocuments []string `json:"specific_documents"`
}

type InitializationStatusRequest struct {
	JobID string `json:"job_id"`
}

type ExpertSearchRequest struct {
	Query            string `json:"query"`
	FilterDepartment string `json:"filter_department"`
	MaxResults       int32  `json:"max_results"`
}

type DepartmentSearchRequest struct {
	Query      string `json:"query"`
	MaxResults int32  `json:"max_results"`
}

type TopicAnalysisRequest struct {
	Query string `json:"query"`
}

type ArticleSearchRequest struct {
	Query      string `json:"query"`
	MaxResults int32  `json:"max_results"`
}

// RESPONSE STRUCTURES
type InitializeResponse struct {
	JobID                     string `json:"job_id"`
	Status                    string `json:"status"`
	TotalArticles             int32  `json:"total_articles"`
	ProcessedArticles         int32  `json:"processed_articles"`
	EstimatedRemainingSeconds int32  `json:"estimated_remaining_seconds"`
	StartedAt                 string `json:"started_at"`
	CurrentArticle            string `json:"current_article"`
}

type InitializationStatusResponse struct {
	JobID                   string           `json:"job_id"`
	Status                  string           `json:"status"`
	ProcessedArticles       int32            `json:"processed_articles"`
	TotalArticles           int32            `json:"total_articles"`
	ProgressPercentage      float32          `json:"progress_percentage"`
	CurrentArticle          string           `json:"current_article"`
	EstimatedCompletionTime string           `json:"estimated_completion_time"`
	Errors                  []string         `json:"errors"`
	TopicsStatistics        map[string]int32 `json:"topics_statistics"`
}

type ExpertSearchResponse struct {
	Experts    []*Expert `json:"experts"`
	TotalFound int32     `json:"total_found"`
}

type Expert struct {
	AuthorID         string  `json:"author_id"`
	NameRu           string  `json:"name_ru"`
	Department       string  `json:"department"`
	ExpertiseScore   float32 `json:"expertise_score"`
	ArticleCount     int32   `json:"article_count"`
	TotalCitations   int32   `json:"total_citations"`
	LastActivityYear int32   `json:"last_activity_year"`
}

type DepartmentSearchResponse struct {
	Departments []*Department `json:"departments"`
	TotalFound  int32         `json:"total_found"`
}

type Department struct {
	OrganizationID string   `json:"organization_id"`
	Name           string   `json:"name"`
	StrengthScore  float32  `json:"strength_score"`
	ExpertCount    int32    `json:"expert_count"`
	TotalArticles  int32    `json:"total_articles"`
	KeyAuthors     []string `json:"key_authors"`
}

type TopicAnalysisResponse struct {
	MainTopic        string   `json:"main_topic"`
	RelatedTopics    []string `json:"related_topics"`
	TopicDescription string   `json:"topic_description"`
}

type ArticleSearchResponse struct {
	Articles   []*Article `json:"articles"`
	TotalFound int32      `json:"total_found"`
}

type Article struct {
	DocumentID     string   `json:"document_id"`
	TitleRu        string   `json:"title_ru"`
	AbstractRu     string   `json:"abstract_ru"`
	Year           int32    `json:"year"`
	Authors        []string `json:"authors"`
	RelevanceScore float32  `json:"relevance_score"`
	MatchedTopics  []string `json:"matched_topics"`
}

type ArticleTopic struct {
	TopicName  string  `json:"topic_name"`
	Confidence float32 `json:"confidence"`
	TopicType  string  `json:"topic_type"`
}

// Структура для передачи статьи с векторами в семантический поиск
type SemanticArticle struct {
	DocumentID        string
	TitleRu           string
	AbstractRu        string
	TitleEmbedding    []byte // marshaled []float32
	AbstractEmbedding []byte // marshaled []float32
}

type SemanticArticleSearchRequest struct {
	Query      string
	MaxResults int
	Offset     int
}

type SemanticArticleResult struct {
	DocumentID     string
	TitleRu        string
	AbstractRu     string
	Year           int32
	Authors        []string
	RelevanceScore float32
	MatchedTopics  []string
}

type SemanticArticleSearchResponse struct {
	Results    []SemanticArticleResult
	TotalFound int
}
type SemanticSearchRequest struct {
	QueryVector []float32              `json:"query_vector"`
	Articles    []SemanticArticleForML `json:"articles"`
	MaxResults  int                    `json:"max_results"`
}

type SemanticArticleForML struct {
	DocumentID        string    `json:"document_id"`
	TitleRu           string    `json:"title_ru"`
	AbstractRu        string    `json:"abstract_ru"`
	TitleEmbedding    []float32 `json:"title_embedding"`
	AbstractEmbedding []float32 `json:"abstract_embedding"`
}

type SemanticSearchResponse struct {
	Results    []SemanticArticleMLResult `json:"results"`
	TotalFound int                       `json:"total_found"`
}

type SemanticArticleMLResult struct {
	DocumentID      string   `json:"document_id"`
	RelevanceScore  float32  `json:"relevance_score"`
	MatchedConcepts []string `json:"matched_concepts"`
}

type QueryAnalysisRequest struct {
	UserQuery string `json:"user_query"`
	Context   string `json:"context"`
}

type QueryAnalysisResponse struct {
	QueryVector []float32 `json:"query_vector"`
}

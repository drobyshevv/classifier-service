package ml

import (
	"context"
	"fmt"

	agentv1 "github.com/drobyshevv/proto-ai-agent/gen/go/proto/ai_agent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MLServiceClient struct {
	client agentv1.AIAnalysisServiceClient
	conn   *grpc.ClientConn
}

func NewMLServiceClient(grpcAddr string) (*MLServiceClient, error) {
	conn, err := grpc.NewClient(grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AI agent: %w", err)
	}

	return &MLServiceClient{
		client: agentv1.NewAIAnalysisServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *MLServiceClient) AnalyzeUserQuery(
	ctx context.Context,
	userQuery string,
	contextType string,
) (*agentv1.QueryAnalysisResponse, error) {

	resp, err := c.client.AnalyzeUserQuery(ctx, &agentv1.QueryAnalysisRequest{
		UserQuery: userQuery,
		Context:   contextType,
	})
	if err != nil {
		return nil, fmt.Errorf("AnalyzeUserQuery failed: %w", err)
	}

	return resp, nil
}

func (c *MLServiceClient) SemanticArticleSearch(
	ctx context.Context,
	req *agentv1.SemanticSearchRequest,
) (*agentv1.SemanticSearchResponse, error) {

	resp, err := c.client.SemanticArticleSearch(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("SemanticArticleSearch failed: %w", err)
	}

	return resp, nil
}

func (c *MLServiceClient) AnalyzeExpertsByTopic(
	ctx context.Context,
	req *agentv1.ExpertAnalysisRequest,
) (*agentv1.ExpertAnalysisResponse, error) {

	resp, err := c.client.AnalyzeExpertsByTopic(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("AnalyzeExpertsByTopic failed: %w", err)
	}

	return resp, nil
}

func (c *MLServiceClient) AnalyzeDepartmentsByTopic(
	ctx context.Context,
	req *agentv1.DepartmentAnalysisRequest,
) (*agentv1.DepartmentAnalysisResponse, error) {

	resp, err := c.client.AnalyzeDepartmentsByTopic(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("AnalyzeDepartmentsByTopic failed: %w", err)
	}

	return resp, nil
}

func (c *MLServiceClient) AnalyzeArticleTopics(ctx context.Context, documentID, title, abstract string) (*ArticleAnalysisResult, error) {
	resp, err := c.client.AnalyzeArticleTopics(ctx, &agentv1.ArticleAnalysisRequest{
		DocumentId: documentID,
		TitleRu:    title,
		AbstractRu: abstract,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to analyze article topics: %w", err)
	}

	result := &ArticleAnalysisResult{
		Topics:            convertProtoTopics(resp.Topics),
		TitleEmbedding:    resp.TitleEmbedding,
		AbstractEmbedding: resp.AbstractEmbedding,
	}

	return result, nil
}

func (c *MLServiceClient) Close() {
	c.conn.Close()
}

type ArticleAnalysisResult struct {
	Topics            []ArticleTopic
	TitleEmbedding    []byte
	AbstractEmbedding []byte
}

type ArticleTopic struct {
	TopicName  string
	Confidence float32
	TopicType  string
}

func convertProtoTopics(protoTopics []*agentv1.ArticleTopic) []ArticleTopic {
	topics := make([]ArticleTopic, len(protoTopics))
	for i, topic := range protoTopics {
		topics[i] = ArticleTopic{
			TopicName:  topic.TopicName,
			Confidence: topic.Confidence,
			TopicType:  topic.TopicType,
		}
	}
	return topics
}

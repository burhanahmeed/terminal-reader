package embed

import (
	"context"
	"errors"
	"fmt"
	"os"

	genai "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiEmbedder struct {
	client *genai.Client
	model  *genai.EmbeddingModel
}

func NewGeminiEmbedder() (*GeminiEmbedder, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GEMINI_API_KEY not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.EmbeddingModel("text-embedding-004")
	return &GeminiEmbedder{client: client, model: model}, nil
}

func (g *GeminiEmbedder) EmbedText(text string) ([]float32, error) {
	if g.model == nil {
		return nil, errors.New("GeminiEmbedder: model not initialized")
	}

	ctx := context.Background()
	resp, err := g.model.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}

	if len(resp.Embedding.Values) == 0 {
		return nil, errors.New("no embedding returned")
	}
	return resp.Embedding.Values, nil
}

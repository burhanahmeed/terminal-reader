package llm

import (
	"context"
	"os"

	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	model *genai.GenerativeModel
}

func NewGeminiClient() (*GeminiClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}
	return &GeminiClient{model: client.GenerativeModel("gemini-2.5-flash")}, nil
}

func (g *GeminiClient) Generate(prompt string) (string, error) {
	resp, err := g.model.GenerateContent(context.Background(), genai.Text(prompt))
	if err != nil {
		return "", err
	}
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response")
	}
	return fmt.Sprint(resp.Candidates[0].Content.Parts[0]), nil
}

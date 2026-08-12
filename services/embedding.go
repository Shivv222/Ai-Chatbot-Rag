package services

import (
	"context"
	"os"

	"google.golang.org/genai"
)

func GenerateEmbedding(text string) ([]float32, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	dim := int32(768)

	result, err := client.Models.EmbedContent(
		ctx,
		"gemini-embedding-001",
		[]*genai.Content{
			genai.NewContentFromText(text, genai.RoleUser),
		},
		&genai.EmbedContentConfig{
			TaskType:             "RETRIEVAL_DOCUMENT",
			OutputDimensionality: &dim,
		},
	)

	if err != nil {
		return nil, err
	}

	return result.Embeddings[0].Values, nil
}

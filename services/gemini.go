package services

import (
	"context"
	"os"

	"google.golang.org/genai"
)

func AskGemini(prompt string) (string, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return "", err
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-3.6-flash",
		genai.Text(prompt),
		nil,
	)

	if err != nil {
		return "", err
	}

	return result.Text(), nil
}

func GenerateChatTitle(firstMessage string) (string, error) {

	prompt := `Generate a short chat title (maximum 5 words).
Return ONLY the title.
Do not use quotes.
Do not explain anything.

User message:
` + firstMessage

	title, err := AskGemini(prompt)
	if err != nil {
		return "", err
	}

	return title, nil
}

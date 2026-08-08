package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func main() {
	_ = godotenv.Load()

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		panic(err)
	}

	pager, err := client.Models.List(ctx, nil)
	if err != nil {
		panic(err)
	}

	for pager.Next() {
		model := pager.Cur()

		fmt.Println(model.Name)

		for _, action := range model.SupportedActions {
			fmt.Println("  -", action)
		}
		fmt.Println("----------------")
	}

	if err := pager.Err(); err != nil {
		panic(err)
	}
}

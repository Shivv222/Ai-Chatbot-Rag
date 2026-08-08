package main

import (
	"fmt"

	"Ai-Chatbot-Rag/services"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	embedding, err := services.GenerateEmbedding("Hello World")
	if err != nil {
		panic(err)
	}

	fmt.Println("Embedding Length:", len(embedding))
}
package main

import (
	"fmt"
	"log"

	"Ai-Chatbot-Rag/config"
	"Ai-Chatbot-Rag/services"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Connect to Docker PostgreSQL
	config.ConnectDB()

	question := "What is the education qualification?"

	answer, sources, err := services.AskRAG(question, 1)

	if err != nil {
		log.Fatal("RAG Error:", err)
	}

	fmt.Println("\n===== RAG ANSWER =====")
	fmt.Println(answer)

	fmt.Println("\n===== SOURCES =====")

	for _, source := range sources {
		fmt.Printf(
			"Document ID: %d | File: %s | Similarity: %.4f\n",
			source.DocumentID,
			source.FileName,
			source.Similarity,
		)
	}
}

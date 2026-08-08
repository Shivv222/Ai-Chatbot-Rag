package main

import (
	"fmt"
	"log"

	"Ai-Chatbot-Rag/config"
	"Ai-Chatbot-Rag/repository"
	"Ai-Chatbot-Rag/services"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	config.ConnectDB()

	query := "What is the education qualification?"

	embedding, err := services.GenerateEmbedding(query)
	if err != nil {
		log.Fatal(err)
	}

	chunks, err := repository.SearchSimilarChunks(
		embedding,
		3,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Top Similar Chunks:")

	for i, chunk := range chunks {
		fmt.Printf("\n--- Chunk %d ---\n", i+1)
		fmt.Println(chunk)
	}
}

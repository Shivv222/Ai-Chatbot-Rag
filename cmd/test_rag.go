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

	answer, err := services.AskRAG(question)
	if err != nil {
		log.Fatal("RAG Error:", err)
	}

	fmt.Println("\n===== RAG ANSWER =====")
	fmt.Println(answer)
}
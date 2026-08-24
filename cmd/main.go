package main

import (
	"Ai-Chatbot-Rag/config"
	"Ai-Chatbot-Rag/controllers"
	"Ai-Chatbot-Rag/middleware"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	router.POST("/register", controllers.Register)

	router.POST("/login", controllers.Login)

	router.GET("/profile", middleware.AuthMiddleware(), controllers.Profile)

	router.POST("/chat", middleware.AuthMiddleware(), controllers.Chat)

	router.POST("/documents/upload", middleware.AuthMiddleware(), controllers.UploadDocument)

	router.GET("/documents", middleware.AuthMiddleware(), controllers.GetDocuments)

	router.GET("/chat/history", middleware.AuthMiddleware(), controllers.GetChatHistory)

	router.POST("/session", middleware.AuthMiddleware(), controllers.CreateSession)

	router.GET("/sessions", middleware.AuthMiddleware(), controllers.GetSessions)

	router.GET("/sessions/:id", middleware.AuthMiddleware(), controllers.GetSessionChats)

	router.PUT("/sessions/:id", middleware.AuthMiddleware(), controllers.UpdateSession)

	router.DELETE("/sessions/:id", middleware.AuthMiddleware(), controllers.DeleteSession)

	router.Run((":8080"))
}

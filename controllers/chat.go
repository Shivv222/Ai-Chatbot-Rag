package controllers

import (
	"Ai-Chatbot-Rag/repository"
	"Ai-Chatbot-Rag/services"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	SessionID int    `json:"session_id" binding:"required"`
	Prompt    string `json:"prompt" binding:"required"`
}

type CreateSessionRequest struct {
	FirstMessage string `json:"first_message" binding:"required"`
}

func Chat(c *gin.Context) {

	var req ChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Prompt is required",
		})
		return
	}

	userID := c.GetInt("userID")

	// 1. Load recent conversation history
	recentChats, err := repository.GetRecentChats(userID, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load recent chats",
		})
		return
	}

	var builder strings.Builder

	builder.WriteString("Previous Conversation:\n\n")

	for i := len(recentChats) - 1; i >= 0; i-- {
		builder.WriteString(fmt.Sprintf(
			"User: %s\n",
			recentChats[i].UserMessage,
		))

		builder.WriteString(fmt.Sprintf(
			"AI: %s\n\n",
			recentChats[i].AIResponse,
		))
	}

	// 2. Generate RAG answer
	ragResponse, err := services.AskRAG(req.Prompt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate RAG response",
		})
		return
	}

	// 3. Add RAG response to conversation context
	builder.WriteString("Document Context Answer:\n")
	builder.WriteString(ragResponse)
	builder.WriteString("\n\n")

	builder.WriteString(fmt.Sprintf(
		"Current User: %s",
		req.Prompt,
	))

	fullPrompt := builder.String()

	// 4. Generate final response using Gemini
	response, err := services.AskGemini(fullPrompt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 5. Save chat history
	err = repository.SaveChat(
		req.SessionID,
		userID,
		req.Prompt,
		response,
	)

	if err != nil {
		fmt.Println("SAVE CHAT ERROR:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 6. Return response
	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}

func GetChatHistory(c *gin.Context) {

	userID := c.GetInt("userID")

	chats, err := repository.GetChatHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch chat history",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history": chats,
	})
}

func CreateSession(c *gin.Context) {

	var req CreateSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title is required",
		})
		return
	}

	userID := c.GetInt("userID")

	title, err := services.GenerateChatTitle(req.FirstMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate chat title",
		})
		return
	}

	sessionID, err := repository.CreateSession(
		userID,
		title,
	)

	if err != nil {
		fmt.Println("CREATE SESSION ERROR:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Session created successfully",
		"session_id": sessionID,
	})
}

func GetSessions(c *gin.Context) {

	userID := c.GetInt("userID")

	sessions, err := repository.GetSessions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch sessions",
		})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

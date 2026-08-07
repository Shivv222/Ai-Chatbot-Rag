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
		builder.WriteString(fmt.Sprintf("User: %s\n", recentChats[i].UserMessage))
		builder.WriteString(fmt.Sprintf("AI: %s\n\n", recentChats[i].AIResponse))
	}

	builder.WriteString(fmt.Sprintf("Current User: %s", req.Prompt))

	fullPrompt := builder.String()

	response, err := services.AskGemini(fullPrompt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = repository.SaveChat(
		req.SessionID,
		userID,
		req.Prompt,
		response,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save chat history",
		})
		return
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create session",
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

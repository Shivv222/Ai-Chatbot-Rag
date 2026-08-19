package controllers

import (
	"Ai-Chatbot-Rag/repository"
	"Ai-Chatbot-Rag/services"
	"fmt"
	"net/http"

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
			"error": "session_id and prompt are required",
		})
		return
	}

	userID := c.GetInt("userID")

	// 1. Verify that this session belongs to the logged-in user
	session, err := repository.GetSession(req.SessionID, userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Session not found",
		})
		return
	}

	_ = session

	// 2. Load recent chats from THIS session only
	recentChats, err := repository.GetRecentChats(
		userID,
		req.SessionID,
		5,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load recent chats",
		})
		return
	}

	// 3. RAG + ONE Gemini call
	response, sources, err := services.AskRAG(
		req.Prompt,
		userID,
		recentChats,
	)

	if err != nil {

		fmt.Println("RAG ERROR:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 4. Save chat history
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

	// 5. Return response and sources
	c.JSON(http.StatusOK, gin.H{
		"response": response,
		"sources":  sources,
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

	title := req.FirstMessage

	if len(title) > 50 {
		title = title[:50] + "..."
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

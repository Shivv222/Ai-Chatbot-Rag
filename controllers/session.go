package controllers

import (
	"net/http"
	"strconv"

	"Ai-Chatbot-Rag/repository"

	"github.com/gin-gonic/gin"
)

type UpdateSessionRequest struct {
	Title string `json:"title" binding:"required"`
}

func GetSessionChats(c *gin.Context) {

	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid session ID",
		})
		return
	}

	userID := c.GetInt("userID")

	chats, err := repository.GetSessionChats(sessionID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load session chats",
		})
		return
	}

	c.JSON(http.StatusOK, chats)
}

func UpdateSession(c *gin.Context) {

	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid session ID",
		})
		return
	}

	var req UpdateSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title is required",
		})
		return
	}

	userID := c.GetInt("userID")

	err = repository.UpdateSessionTitle(
		sessionID,
		userID,
		req.Title,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Session updated successfully",
	})
}

func DeleteSession(c *gin.Context) {

	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid session ID",
		})
		return
	}

	userID := c.GetInt("userID")

	err = repository.DeleteSession(
		sessionID,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Session deleted successfully",
	})
}

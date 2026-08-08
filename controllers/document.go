package controllers

import (
	"log"
	"net/http"
	"path/filepath"

	"Ai-Chatbot-Rag/repository"
	"Ai-Chatbot-Rag/services"

	"github.com/gin-gonic/gin"
)

func UploadDocument(c *gin.Context) {

	userID := c.GetInt("userID")

	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Document is required",
		})
		return
	}

	path := "uploads/" + filepath.Base(file.Filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save document",
		})
		return
	}

	text, err := services.ExtractTextFromPDF(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	chunks := services.SplitIntoChunks(text, 100)

	println("Number of Chunks:", len(chunks))

	for i, chunk := range chunks {

		println("------------")

		println("Chunk", i+1)

		println(chunk)
	}

	println(text)

	documentID, err := repository.SaveDocument(
		userID,
		file.Filename,
		path,
	)

	if err != nil {
		log.Println("SAVE DOCUMENT ERROR:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	for i, chunk := range chunks {

		chunkID, err := repository.SaveDocumentChunk(
			documentID,
			chunk,
			i,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save chunk",
			})
			return
		}

		embedding, err := services.GenerateEmbedding(chunk)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate embedding",
			})
			return
		}

		err = repository.SaveChunkEmbedding(
			chunkID,
			embedding,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document uploaded successfully",
		"file":    file.Filename,
	})
}

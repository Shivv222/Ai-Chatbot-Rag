package services

import (
	"fmt"
	"strings"

	"Ai-Chatbot-Rag/models"
	"Ai-Chatbot-Rag/repository"
)

type RAGSource struct {
	DocumentID int     `json:"document_id"`
	FileName   string  `json:"file_name"`
	Similarity float64 `json:"similarity"`
}

func AskRAG(
	question string,
	userID int,
	recentChats []models.ChatHistory,
) (string, []RAGSource, error) {

	// 1. Generate embedding for the user's question
	embedding, err := GenerateEmbedding(question)
	if err != nil {
		return "", nil, err
	}

	// 2. Search relevant chunks from this user's documents
	results, err := repository.SearchSimilarChunks(
		embedding,
		3,
		userID,
	)
	if err != nil {
		return "", nil, err
	}

	// 3. Build document context
	var contextBuilder strings.Builder

	for _, result := range results {

		contextBuilder.WriteString(
			fmt.Sprintf(
				"Source Document: %s\n%s\n\n",
				result.FileName,
				result.ChunkText,
			),
		)
	}

	context := contextBuilder.String()

	// 4. Create source list
	var sources []RAGSource

	seen := make(map[int]bool)

	for _, result := range results {

		if !seen[result.DocumentID] {

			sources = append(sources, RAGSource{
				DocumentID: result.DocumentID,
				FileName:   result.FileName,
				Similarity: result.Similarity,
			})

			seen[result.DocumentID] = true
		}
	}

	// 5. Build conversation history
	var historyBuilder strings.Builder

	for i := len(recentChats) - 1; i >= 0; i-- {

		historyBuilder.WriteString(
			fmt.Sprintf(
				"User: %s\nAI: %s\n\n",
				recentChats[i].UserMessage,
				recentChats[i].AIResponse,
			),
		)
	}

	history := historyBuilder.String()

	// 6. Create ONE Gemini prompt
	prompt := fmt.Sprintf(`
You are a helpful AI assistant.

Answer the user's current question using the provided document context.

CONVERSATION HISTORY:
%s

DOCUMENT CONTEXT:
%s

CURRENT USER QUESTION:
%s

Instructions:

- Use the retrieved document context when it is relevant.
- Use the conversation history to understand follow-up questions.
- Do not invent information.
- If the answer is not available in the document context, clearly say so.
- Keep the answer relevant and concise.
`, history, context, question)

	// 7. ONE Gemini call
	answer, err := AskGemini(prompt)
	if err != nil {
		return "", nil, err
	}

	// 8. Return answer + sources
	return answer, sources, nil
}

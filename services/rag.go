package services

import (
	"fmt"
	"strings"

	"Ai-Chatbot-Rag/repository"
)

type RAGSource struct {
	DocumentID int     `json:"document_id"`
	FileName   string  `json:"file_name"`
	Similarity float64 `json:"similarity"`
}

func AskRAG(question string, userID int) (string, []RAGSource, error) {

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

	// 5. Create prompt
	prompt := fmt.Sprintf(`
You are a helpful AI assistant.

Answer the user's question using the provided document context.

DOCUMENT CONTEXT:
%s

USER QUESTION:
%s

Instructions:

- Use the retrieved document context to answer the question.
- Do not invent information.
- Keep the answer relevant to the user's question.
`, context, question)

	// 6. Ask Gemini
	answer, err := AskGemini(prompt)
	if err != nil {
		return "", nil, err
	}

	// 7. Return answer + sources
	return answer, sources, nil
}

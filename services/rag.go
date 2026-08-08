package services

import (
	"fmt"
	"strings"

	"Ai-Chatbot-Rag/repository"
)

func AskRAG(question string) (string, error) {

	// 1. Generate embedding for the user's question
	embedding, err := GenerateEmbedding(question)
	if err != nil {
		return "", err
	}

	// 2. Search the most relevant PDF chunks
	chunks, err := repository.SearchSimilarChunks(embedding, 3)
	if err != nil {
		return "", err
	}

	// 3. Combine retrieved chunks into context
	context := strings.Join(chunks, "\n\n")

	// 4. Create prompt for Gemini
	prompt := fmt.Sprintf(`
You are a helpful AI assistant.

Answer the user's question using the provided document context.

DOCUMENT CONTEXT:
%s

USER QUESTION:
%s

Instructions:
- Answer using the document context.
- If the answer is not present in the document, say that the information is not available in the document.
- Do not invent information.
`, context, question)

	// 5. Ask Gemini
	answer, err := AskGemini(prompt)
	if err != nil {
		return "", err
	}

	return answer, nil
}

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
		5,
		userID,
	)
	if err != nil {
		return "", nil, err
	}

	// Strict grounding threshold
	const similarityThreshold = 0.65

	// Keep only sufficiently relevant chunks
	var relevantResults []repository.SimilarChunk

	for _, result := range results {
		if result.Similarity >= similarityThreshold {
			relevantResults = append(relevantResults, result)
		}
	}

	// No sufficiently relevant document context found
	if len(relevantResults) == 0 {
		return "I couldn't find enough relevant information in the uploaded documents to answer this question.", nil, nil
	}

	// Use only relevant results from this point onward
	results = relevantResults

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
You are a document question-answering assistant.

Your job is to answer the CURRENT USER QUESTION using ONLY the information
contained in the DOCUMENT CONTEXT and relevant CONVERSATION HISTORY.

CONVERSATION HISTORY:
%s

DOCUMENT CONTEXT:
%s

CURRENT USER QUESTION:
%s

STRICT RULES:

1. Use the DOCUMENT CONTEXT as the primary source of truth.

2. Do NOT use outside knowledge, assumptions, guesses, or information
   that is not supported by the DOCUMENT CONTEXT.

3. If the answer is explicitly available in the DOCUMENT CONTEXT,
   answer it clearly and directly.

4. If the CURRENT USER QUESTION is a follow-up question, you may use the
   CONVERSATION HISTORY to understand what the user is referring to.
   However, the actual factual answer must still be supported by the
   DOCUMENT CONTEXT.

5. If the DOCUMENT CONTEXT does not contain enough information to answer
   the question, respond exactly:
   "I could not find this information in the uploaded document."

6. Do NOT try to answer from general knowledge when the document does
   not contain the answer.

7. Do NOT invent names, dates, numbers, percentages, facts, or conclusions.

8. Keep the answer concise and directly related to the user's question.

9. When answering using document information, mention the relevant
   source document filename when it helps clarify where the information
   came from.

10. Do not cite or mention a document unless the information used in
    the answer actually comes from that document.

Answer:
`, history, context, question)

	// 7. ONE Gemini call
	answer, err := AskGemini(prompt)
	if err != nil {
		return "", nil, err
	}

	// 8. Return answer + sources
	return answer, sources, nil
}

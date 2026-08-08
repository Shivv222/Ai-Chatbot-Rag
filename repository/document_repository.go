package repository

import (
	"Ai-Chatbot-Rag/config"
	"fmt"
	"strings"
)

func SaveDocument(userID int, fileName, filePath string) (int, error) {

	query := `
	INSERT INTO documents
	(user_id, file_name, file_path)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	var documentID int

	err := config.DB.QueryRow(
		query,
		userID,
		fileName,
		filePath,
	).Scan(&documentID)

	if err != nil {
		return 0, err
	}

	return documentID, nil
}

func SaveDocumentChunk(documentID int, chunk string, index int) (int, error) {

	query := `
	INSERT INTO document_chunks
	(document_id, chunk_text, chunk_index)
	VALUES ($1,$2,$3)
	RETURNING id
	`

	var chunkID int

	err := config.DB.QueryRow(
		query,
		documentID,
		chunk,
		index,
	).Scan(&chunkID)

	if err != nil {
		return 0, err
	}

	return chunkID, nil
}

func SaveChunkEmbedding(chunkID int, embedding []float32) error {

	values := make([]string, len(embedding))

	for i, v := range embedding {
		values[i] = fmt.Sprintf("%f", v)
	}

	vector := "[" + strings.Join(values, ",") + "]"

	query := `
	UPDATE document_chunks
	SET embedding = $1
	WHERE id = $2
	`

	_, err := config.DB.Exec(
		query,
		vector,
		chunkID,
	)

	return err
}

func SearchSimilarChunks(
	embedding []float32,
	limit int,
) ([]string, error) {

	values := make([]string, len(embedding))

	for i, v := range embedding {
		values[i] = fmt.Sprintf("%f", v)
	}

	vector := "[" + strings.Join(values, ",") + "]"

	query := `
	SELECT chunk_text
	FROM document_chunks
	WHERE embedding IS NOT NULL
	ORDER BY embedding <=> $1
	LIMIT $2
	`

	rows, err := config.DB.Query(
		query,
		vector,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []string

	for rows.Next() {

		var chunk string

		if err := rows.Scan(&chunk); err != nil {
			return nil, err
		}

		chunks = append(chunks, chunk)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chunks, nil
}

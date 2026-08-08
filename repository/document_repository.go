package repository

import (
	"Ai-Chatbot-Rag/config"
	"fmt"
	"strings"
)

type SimilarChunk struct {
	ChunkText  string
	DocumentID int
	FileName   string
	Similarity float64
}

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
	userID int,
) ([]SimilarChunk, error) {

	values := make([]string, len(embedding))

	for i, v := range embedding {
		values[i] = fmt.Sprintf("%f", v)
	}

	vector := "[" + strings.Join(values, ",") + "]"

	query := `
	SELECT
		dc.chunk_text,
		d.id,
		d.file_name,
		1 - (dc.embedding <=> $1) AS similarity
	FROM document_chunks dc
	JOIN documents d
		ON dc.document_id = d.id
	WHERE dc.embedding IS NOT NULL
	AND d.user_id = $3
	ORDER BY dc.embedding <=> $1
	LIMIT $2
	`

	rows, err := config.DB.Query(
		query,
		vector,
		limit,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []SimilarChunk

	for rows.Next() {

		var result SimilarChunk

		if err := rows.Scan(
			&result.ChunkText,
			&result.DocumentID,
			&result.FileName,
			&result.Similarity,
		); err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

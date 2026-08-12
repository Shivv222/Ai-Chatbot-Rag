package services

import "strings"

func SplitIntoChunks(text string, chunkSize int) []string {

	var chunks []string

	words := strings.Fields(text)

	for i := 0; i < len(words); i += chunkSize {

		end := i + chunkSize

		if end > len(words) {
			end = len(words)
		}

		chunk := strings.Join(words[i:end], " ")

		chunks = append(chunks, chunk)
	}

	return chunks
}
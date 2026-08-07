package repository

import (
	"Ai-Chatbot-Rag/config"
)

func CreateSession(userID int, title string) (int, error) {

	query := `
	INSERT INTO chat_sessions(user_id, title)
	VALUES($1, $2)
	RETURNING id
	`

	var sessionID int

	err := config.DB.QueryRow(
		query,
		userID,
		title,
	).Scan(&sessionID)

	if err != nil {
		return 0, err
	}

	return sessionID, nil
}
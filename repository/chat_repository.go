package repository

import (
	"Ai-Chatbot-Rag/config"
	"Ai-Chatbot-Rag/models"
)

func SaveChat(sessionID int, userID int, userMessage, aiResponse string) error {

	query := `
	INSERT INTO chat_history
	(session_id, user_id, user_message, ai_response)
	VALUES ($1, $2, $3, $4)
	`

	_, err := config.DB.Exec(
		query,
		sessionID,
		userID,
		userMessage,
		aiResponse,
	)

	return err
}

func GetChatHistory(userID int) ([]models.ChatHistory, error) {

	query := `
	SELECT id, user_id, user_message, ai_response, created_at
	FROM chat_history
	WHERE user_id = $1
	ORDER BY created_at ASC
	`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []models.ChatHistory

	for rows.Next() {
		var chat models.ChatHistory

		err := rows.Scan(
			&chat.ID,
			&chat.UserID,
			&chat.UserMessage,
			&chat.AIResponse,
			&chat.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		chats = append(chats, chat)
	}

	return chats, nil
}

func GetRecentChats(userID int, sessionID int, limit int) ([]models.ChatHistory, error) {

	query := `
	SELECT id, user_id, user_message, ai_response, created_at
	FROM chat_history
	WHERE user_id = $1
	AND session_id = $2
	ORDER BY created_at DESC
	LIMIT $3
	`

	rows, err := config.DB.Query(
		query,
		userID,
		sessionID,
		limit,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var chats []models.ChatHistory

	for rows.Next() {

		var chat models.ChatHistory

		err := rows.Scan(
			&chat.ID,
			&chat.UserID,
			&chat.UserMessage,
			&chat.AIResponse,
			&chat.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

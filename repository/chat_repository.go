package repository

import (
	"Ai-Chatbot-Rag/config"
	"Ai-Chatbot-Rag/models"
)

func SaveChat(userID int, userMessage, aiResponse string) error {

	query := `
	INSERT INTO chat_history
	(user_id, user_message, ai_response)
	VALUES ($1, $2, $3)
	`

	_, err := config.DB.Exec(
		query,
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

func GetRecentChats(userID int, limit int) ([]models.ChatHistory, error) {

	query := `
	SELECT id, user_id, user_message, ai_response, created_at
	FROM chat_history
	WHERE user_id = $1
	ORDER BY created_at DESC
	LIMIT $2
	`

	rows, err := config.DB.Query(query, userID, limit)
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
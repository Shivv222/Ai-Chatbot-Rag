package repository

import (
	"Ai-Chatbot-Rag/config"
	"Ai-Chatbot-Rag/models"
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

// When the frontend opens, it should be able to fetch all the user's chat sessions.
func GetSessions(userID int) ([]models.ChatSession, error) {

	query := `
	SELECT id, user_id, title, created_at
	FROM chat_sessions
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.ChatSession

	for rows.Next() {
		var session models.ChatSession

		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.Title,
			&session.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func GetSessionChats(sessionID int, userID int) ([]models.ChatHistory, error) {

	query := `
	SELECT
		id,
		session_id,
		user_id,
		user_message,
		ai_response,
		created_at
	FROM chat_history
	WHERE session_id = $1
	AND user_id = $2
	ORDER BY created_at ASC
	`

	rows, err := config.DB.Query(query, sessionID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []models.ChatHistory

	for rows.Next() {

		var chat models.ChatHistory

		err := rows.Scan(
			&chat.ID,
			&chat.SessionID,
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

func UpdateSessionTitle(sessionID int, userID int, title string) error {

	query := `
	UPDATE chat_sessions
	SET title = $1
	WHERE id = $2
	AND user_id = $3
	`

	_, err := config.DB.Exec(
		query,
		title,
		sessionID,
		userID,
	)

	return err
}

func DeleteSession(sessionID int, userID int) error {

	query := `
	DELETE FROM chat_sessions
	WHERE id = $1
	AND user_id = $2
	`

	_, err := config.DB.Exec(
		query,
		sessionID,
		userID,
	)

	return err
}
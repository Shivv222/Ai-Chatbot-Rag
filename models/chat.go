package models

import "time"

type ChatHistory struct {
	ID          int
	SessionID   int
	UserID      int
	UserMessage string
	AIResponse  string
	CreatedAt   time.Time
}

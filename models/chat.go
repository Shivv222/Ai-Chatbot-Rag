package models

import "time"

type ChatHistory struct {
	ID          int
	UserID      int
	UserMessage string
	AIResponse  string
	CreatedAt   time.Time
}

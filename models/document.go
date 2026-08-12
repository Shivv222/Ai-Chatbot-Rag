package models

import "time"

type Document struct {
	ID        int
	UserID    int
	FileName  string
	FilePath  string
	CreatedAt time.Time
}
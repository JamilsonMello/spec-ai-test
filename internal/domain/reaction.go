package domain

import (
	"time"
)

type Reaction struct {
	ID        string    `json:"id"`
	CommentID string    `json:"commentId"`
	UserID    string    `json:"userId"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}

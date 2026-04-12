package domain

import (
	"time"
)

type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"postId"`
	UserID    string    `json:"userId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

func (c *Comment) IsValidContent() bool {
	return len(c.Content) > 0
}

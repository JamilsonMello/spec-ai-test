package domain

import (
	"html"
	"time"
)

type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func NewComment(postID, userID, content string) *Comment {
	return &Comment{
		PostID:  postID,
		UserID:  userID,
		Content: html.EscapeString(content),
	}
}

func (c *Comment) IsValidContent() bool {
	return len(c.Content) > 0 && len(c.Content) <= 500
}

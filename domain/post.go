package domain

import (
	"time"
)

type Post struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (post *Post) IsValidContent() bool {
	return len(post.Content) > 0 && len(post.Content) <= 600
}

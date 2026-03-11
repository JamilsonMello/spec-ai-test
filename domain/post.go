package domain

import (
	"time"
)

// Post represents a text-only post created by a user.
type Post struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// IsValidContent checks if the post content is valid (max 5000 characters).
func (p *Post) IsValidContent() bool {
	return len(p.Content) > 0 && len(p.Content) <= 5000
}

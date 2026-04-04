package domain

import (
	"time"
)

type Post struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"authorId"`
	VideoURL  string    `json:"videoUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (p *Post) IsValidContent() bool {
	return len(p.Content) > 0 && len(p.Content) <= 5000
}

package domain

import "time"

const MaxCommentLength = 400

// Comment represents a user's comment on a post.
type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"postId"`
	AuthorID  string    `json:"authorId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

// IsValidContent checks if the comment content is valid (1–400 characters).
func (c *Comment) IsValidContent() bool {
	return len(c.Content) > 0 && len(c.Content) <= MaxCommentLength
}

package domain

import "time"

// LikeTargetType identifies whether a like targets a post or a comment.
type LikeTargetType string

const (
	LikeTargetPost    LikeTargetType = "POST"
	LikeTargetComment LikeTargetType = "COMMENT"
)

// Like represents a user's like on a post or comment.
type Like struct {
	UserID     string         `json:"userId"`
	TargetID   string         `json:"targetId"`
	TargetType LikeTargetType `json:"targetType"`
	CreatedAt  time.Time      `json:"createdAt"`
}

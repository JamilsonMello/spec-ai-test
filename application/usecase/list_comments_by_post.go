package usecase

import (
	"time"

	"github.com/example/cadastro-de-usuarios/domain"
)

// CommentItem represents a single comment in the list response.
type CommentItem struct {
	ID        string `json:"id"`
	PostID    string `json:"postId"`
	AuthorID  string `json:"authorId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	LikeCount int    `json:"likeCount"`
}

// ListCommentsResponse represents the output for listing comments on a post.
type ListCommentsResponse struct {
	PostLikeCount int           `json:"postLikeCount"`
	Comments      []CommentItem `json:"comments"`
}

// ListCommentsByPostUseCase handles retrieval of comments and like counts for a post.
type ListCommentsByPostUseCase struct {
	PostRepository    domain.PostRepository
	CommentRepository domain.CommentRepository
	LikeRepository    domain.LikeRepository
}

// NewListCommentsByPostUseCase creates a new ListCommentsByPostUseCase.
func NewListCommentsByPostUseCase(
	postRepo domain.PostRepository,
	commentRepo domain.CommentRepository,
	likeRepo domain.LikeRepository,
) *ListCommentsByPostUseCase {
	return &ListCommentsByPostUseCase{
		PostRepository:    postRepo,
		CommentRepository: commentRepo,
		LikeRepository:    likeRepo,
	}
}

// Execute returns all comments for a post along with like counts.
func (uc *ListCommentsByPostUseCase) Execute(postID string) (*ListCommentsResponse, error) {
	_, err := uc.PostRepository.GetPostByID(postID)
	if err != nil {
		return nil, domain.ErrPostNotFound
	}

	postLikeCount, err := uc.LikeRepository.CountLikesByTarget(postID)
	if err != nil {
		return nil, err
	}

	comments, err := uc.CommentRepository.ListCommentsByPostID(postID)
	if err != nil {
		return nil, err
	}

	items := make([]CommentItem, 0, len(comments))
	for _, c := range comments {
		likeCount, err := uc.LikeRepository.CountLikesByTarget(c.ID)
		if err != nil {
			return nil, err
		}
		items = append(items, CommentItem{
			ID:        c.ID,
			PostID:    c.PostID,
			AuthorID:  c.AuthorID,
			Content:   c.Content,
			CreatedAt: c.CreatedAt.Format(time.RFC3339),
			LikeCount: likeCount,
		})
	}

	return &ListCommentsResponse{
		PostLikeCount: postLikeCount,
		Comments:      items,
	}, nil
}

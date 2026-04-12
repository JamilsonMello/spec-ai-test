package usecase

import (
	"time"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type ListCommentsRequest struct {
	PostID string
}

type CommentItem struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ListCommentsResponse struct {
	Comments []CommentItem `json:"comments"`
}

type ListCommentsUseCase struct {
	CommentRepository domain.CommentRepository
}

func NewListCommentsUseCase(repo domain.CommentRepository) *ListCommentsUseCase {
	return &ListCommentsUseCase{
		CommentRepository: repo,
	}
}

func (uc *ListCommentsUseCase) Execute(req ListCommentsRequest) (*ListCommentsResponse, error) {
	comments, err := uc.CommentRepository.GetCommentsByPostID(req.PostID)
	if err != nil {
		return nil, err
	}

	items := make([]CommentItem, 0, len(comments))
	for _, c := range comments {
		items = append(items, CommentItem{
			ID:        c.ID,
			Content:   c.Content,
			UserID:    c.UserID,
			CreatedAt: c.CreatedAt,
		})
	}

	return &ListCommentsResponse{Comments: items}, nil
}

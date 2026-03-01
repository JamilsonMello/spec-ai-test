package usecase

import (
	"errors"
	"time"

	"github.com/example/cadastro-de-usuarios/domain"
)

var ErrLikeUnauthorized = errors.New("usuário não autenticado")

// ToggleLikeRequest represents the input data for toggling a like.
type ToggleLikeRequest struct {
	UserID     string
	TargetID   string
	TargetType domain.LikeTargetType
}

// ToggleLikeResponse represents the output after toggling a like.
type ToggleLikeResponse struct {
	Liked     bool `json:"liked"`
	LikeCount int  `json:"likeCount"`
}

// ToggleLikeUseCase handles the business logic for toggling a like on a post or comment.
type ToggleLikeUseCase struct {
	PostRepository    domain.PostRepository
	CommentRepository domain.CommentRepository
	LikeRepository    domain.LikeRepository
}

// NewToggleLikeUseCase creates a new ToggleLikeUseCase.
func NewToggleLikeUseCase(
	postRepo domain.PostRepository,
	commentRepo domain.CommentRepository,
	likeRepo domain.LikeRepository,
) *ToggleLikeUseCase {
	return &ToggleLikeUseCase{
		PostRepository:    postRepo,
		CommentRepository: commentRepo,
		LikeRepository:    likeRepo,
	}
}

// Execute toggles a like on the given target and returns the new state.
func (uc *ToggleLikeUseCase) Execute(req ToggleLikeRequest) (*ToggleLikeResponse, error) {
	if req.UserID == "" {
		return nil, ErrLikeUnauthorized
	}

	if err := uc.validateTargetExists(req.TargetID, req.TargetType); err != nil {
		return nil, err
	}

	existing, err := uc.LikeRepository.GetLike(req.UserID, req.TargetID)
	if err != nil {
		return nil, err
	}

	liked := false
	if existing != nil {
		if err := uc.LikeRepository.DeleteLike(req.UserID, req.TargetID); err != nil {
			return nil, err
		}
	} else {
		like := &domain.Like{
			UserID:     req.UserID,
			TargetID:   req.TargetID,
			TargetType: req.TargetType,
			CreatedAt:  time.Now(),
		}
		if err := uc.LikeRepository.SaveLike(like); err != nil {
			return nil, err
		}
		liked = true
	}

	count, err := uc.LikeRepository.CountLikesByTarget(req.TargetID)
	if err != nil {
		return nil, err
	}

	return &ToggleLikeResponse{Liked: liked, LikeCount: count}, nil
}

func (uc *ToggleLikeUseCase) validateTargetExists(targetID string, targetType domain.LikeTargetType) error {
	switch targetType {
	case domain.LikeTargetPost:
		_, err := uc.PostRepository.GetPostByID(targetID)
		return err
	case domain.LikeTargetComment:
		_, err := uc.CommentRepository.GetCommentByID(targetID)
		return err
	default:
		return errors.New("tipo de alvo inválido")
	}
}

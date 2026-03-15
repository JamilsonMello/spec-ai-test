package usecase

import (
	"strings"
	"time"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type ListUsersRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	BirthDate string    `json:"birthDate"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type ListUsersResponse struct {
	Users      []UserResponse `json:"users"`
	TotalCount int            `json:"totalCount"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
}

type ListUsersUseCase struct {
	UserRepository domain.UserRepository
}

func NewListUsersUseCase(repo domain.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{
		UserRepository: repo,
	}
}

func (uc *ListUsersUseCase) Execute(req ListUsersRequest) (*ListUsersResponse, error) {
	uc.validateRequest(&req)

	filter := uc.buildFilter(req)

	users, totalCount, err := uc.UserRepository.ListUsers(filter, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	userResponses := uc.mapUsersToResponses(users)

	return uc.buildResponse(userResponses, totalCount, req.Page, req.Limit), nil
}

func (uc *ListUsersUseCase) validateRequest(req *ListUsersRequest) {
	if req.Limit <= 0 || req.Limit > 30 {
		req.Limit = 30
	}

	if req.Page <= 0 {
		req.Page = 1
	}
}

func (uc *ListUsersUseCase) buildFilter(req ListUsersRequest) domain.UserFilter {
	return domain.UserFilter{
		Name:  strings.ToLower(strings.TrimSpace(req.Name)),
		Email: strings.ToLower(strings.TrimSpace(req.Email)),
	}
}

func (uc *ListUsersUseCase) mapUsersToResponses(users []*domain.User) []UserResponse {
	userResponses := make([]UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Surname:   user.Surname,
			Email:     user.Email,
			BirthDate: user.BirthDate.Format("2006-01-02"),
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		})
	}
	return userResponses
}

func (uc *ListUsersUseCase) buildResponse(userResponses []UserResponse, totalCount int, page int, limit int) *ListUsersResponse {
	return &ListUsersResponse{
		Users:      userResponses,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
	}
}

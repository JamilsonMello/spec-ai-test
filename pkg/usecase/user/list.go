package user

import (
	"strings"
	"time"

	"github.com/example/cadastro-de-usuarios/pkg/domain/repository"
)

// ListUsersRequest is the DTO for listing users with filters and pagination.
type ListUsersRequest struct {
	Name  string `json:"name"`  // Optional filter by name
	Email string `json:"email"` // Optional filter by email
	Page  int    `json:"page"`  // Page number (starts from 1)
	Limit int    `json:"limit"` // Items per page (default 30, max 30)
}

// UserResponse is the DTO for user data in list responses.
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	BirthDate string    `json:"birthDate"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

// ListUsersResponse is the DTO for listing users output.
type ListUsersResponse struct {
	Users      []UserResponse `json:"users"`
	TotalCount int            `json:"totalCount"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
}

// ListUsersUseCase handles the business logic for listing users.
type ListUsersUseCase struct {
	UserRepository repository.UserRepository
}

// NewListUsersUseCase creates a new ListUsersUseCase.
func NewListUsersUseCase(repo repository.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{
		UserRepository: repo,
	}
}

// Execute performs the user listing process with filters and pagination.
func (uc *ListUsersUseCase) Execute(req ListUsersRequest) (*ListUsersResponse, error) {
	// Set default limit to 30 if not provided or exceeds max
	if req.Limit <= 0 || req.Limit > 30 {
		req.Limit = 30
	}

	// Set default page to 1 if not provided
	if req.Page <= 0 {
		req.Page = 1
	}

	// Build filter criteria
	filter := repository.UserFilter{
		Name:  strings.ToLower(strings.TrimSpace(req.Name)),
		Email: strings.ToLower(strings.TrimSpace(req.Email)),
	}

	// Retrieve users from repository with filters
	users, totalCount, err := uc.UserRepository.ListUsers(filter, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert domain users to response DTOs (excluding sensitive fields)
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

	return &ListUsersResponse{
		Users:      userResponses,
		TotalCount: totalCount,
		Page:       req.Page,
		Limit:      req.Limit,
	}, nil
}

package usecase

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/example/cadastro-de-usuarios/domain"
)

// Define roles
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// ListUsersRequest defines the request structure for listing users.
type ListUsersRequest struct {
	Name  string
	Email string
	Page  int
	Limit int
	Role  string // Role of the authenticated user making the request
}

// ListUsersResponse defines the response structure for listing users.
type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// UserResponse defines the structure of a user in the response, omitting sensitive fields.
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birthDate"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

// ListUsersUseCase handles the business logic for listing users.
type ListUsersUseCase struct {
	UserRepository UserRepository
}

// NewListUsersUseCase creates a new ListUsersUseCase.
func NewListUsersUseCase(repo UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{
		UserRepository: repo,
	}
}

// Execute performs the user listing process.
func (uc *ListUsersUseCase) Execute(req ListUsersRequest) (*ListUsersResponse, error) {
	// Authorization check
	if req.Role != RoleAdmin {
		return nil, errors.New("access denied: only administrators can list users")
	}

	// Set default pagination values
	if req.Limit <= 0 || req.Limit > 30 {
		req.Limit = 30
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	allUsers, err := uc.UserRepository.ListUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	// Apply filters
	filteredUsers := []domain.User{}
	for _, user := range allUsers {
		nameMatch := true
		emailMatch := true

		if req.Name != "" {
			nameMatch = strings.Contains(strings.ToLower(user.Name), strings.ToLower(req.Name)) ||
				strings.Contains(strings.ToLower(user.Surname), strings.ToLower(req.Name))
		}
		if req.Email != "" {
			emailMatch = strings.Contains(strings.ToLower(user.Email), strings.ToLower(req.Email))
		}

		if nameMatch && emailMatch {
			filteredUsers = append(filteredUsers, user)
		}
	}

	// Sort by CreatedAt in descending order (FR-004)
	sort.Slice(filteredUsers, func(i, j int) bool {
		return filteredUsers[j].CreatedAt.Before(filteredUsers[i].CreatedAt)
	})

	totalUsers := len(filteredUsers)

	// Apply pagination
	start := (req.Page - 1) * req.Limit
	end := start + req.Limit

	if start > totalUsers {
		return &ListUsersResponse{
			Users: []UserResponse{},
			Total: totalUsers,
			Page:  req.Page,
			Limit: req.Limit,
		}, nil
	}
	if end > totalUsers {
		end = totalUsers
	}

	paginatedUsers := filteredUsers[start:end]

	// Format response, omitting sensitive fields
	userResponses := make([]UserResponse, len(paginatedUsers))
	for i, user := range paginatedUsers {
		userResponses[i] = UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Surname:   user.Surname,
			Email:     user.Email,
			BirthDate: user.BirthDate,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		}
	}

	return &ListUsersResponse{
		Users: userResponses,
		Total: totalUsers,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

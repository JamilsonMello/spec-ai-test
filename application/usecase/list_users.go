package usecase

import (
	"strings"
	"time"

	"github.com/example/cadastro-de-usuarios/domain"
)

type ListUsersInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

type UserOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	BirthDate string    `json:"birthDate"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type ListUsersOutput struct {
	Users      []UserOutput `json:"users"`
	TotalCount int          `json:"totalCount"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
}

type ListUsersUseCase struct {
	UserRepository domain.UserRepository
}

func NewListUsersUseCase(repo domain.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{
		UserRepository: repo,
	}
}

func (uc *ListUsersUseCase) validatePagination(limit int, page int) (int, int) {
	if limit <= 0 || limit > 30 {
		limit = 30
	}
	if page <= 0 {
		page = 1
	}
	return limit, page
}

func (uc *ListUsersUseCase) buildFilter(name, email string) domain.UserFilter {
	return domain.UserFilter{
		Name:  strings.ToLower(strings.TrimSpace(name)),
		Email: strings.ToLower(strings.TrimSpace(email)),
	}
}

func (uc *ListUsersUseCase) mapUsersToOutput(users []*domain.User) []UserOutput {
	userResponses := make([]UserOutput, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, UserOutput{
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

func (uc *ListUsersUseCase) Execute(req ListUsersInput) (*ListUsersOutput, error) {

	req.Limit, req.Page = uc.validatePagination(req.Limit, req.Page)

	filter := uc.buildFilter(req.Name, req.Email)

	users, totalCount, repositoryErr := uc.UserRepository.ListUsers(filter, req.Page, req.Limit)
	if repositoryErr != nil {
		return nil, repositoryErr
	}

	userResponses := uc.mapUsersToOutput(users)

	return &ListUsersOutput{
		Users:      userResponses,
		TotalCount: totalCount,
		Page:       req.Page,
		Limit:      req.Limit,
	}, nil
}

package usecase

import "time"

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
func (uc *ListUsersUseCase) Execute(params ListUsersParams) ([]*ListUsersResponse, error) {
	// Set default values for pagination
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 30
	}

	users, err := uc.UserRepository.FindAll(params)
	if err != nil {
		return nil, err
	}

	var response []*ListUsersResponse
	for _, user := range users {
		response = append(response, &ListUsersResponse{
			ID:        user.ID,
			Name:      user.Name,
			Surname:   user.Surname,
			Email:     user.Email,
			BirthDate: user.BirthDate.Format("2006-01-02"),
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		})
	}

	return response, nil
}

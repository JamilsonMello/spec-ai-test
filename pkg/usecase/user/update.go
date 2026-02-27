package user

import (
	"errors"
	"time"

	"github.com/example/cadastro-de-usuarios/pkg/domain/repository"
)

// Custom errors for update user profile validation
var (
	ErrInvalidNameUpdate      = errors.New("nome deve ter entre 2 e 50 caracteres e conter apenas letras e espaços")
	ErrInvalidBirthDateUpdate = errors.New("data de nascimento inválida")
	ErrFutureBirthDateUpdate  = errors.New("data de nascimento não pode ser no futuro")
	ErrUserNotFoundUpdate     = errors.New("usuário não encontrado")
)

// UpdateUserProfileRequest is the DTO for user profile update input.
type UpdateUserProfileRequest struct {
	UserID    string `param:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birthDate"` // YYYY-MM-DD
}

// UpdateUserProfileResponse is the DTO for user profile update output.
type UpdateUserProfileResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	BirthDate string `json:"birthDate"`
}

// UpdateUserProfileUseCase handles the business logic for updating user profile.
type UpdateUserProfileUseCase struct {
	UserRepository repository.UserRepository
}

// NewUpdateUserProfileUseCase creates a new UpdateUserProfileUseCase.
func NewUpdateUserProfileUseCase(repo repository.UserRepository) *UpdateUserProfileUseCase {
	return &UpdateUserProfileUseCase{
		UserRepository: repo,
	}
}

// Execute performs the user profile update process.
func (uc *UpdateUserProfileUseCase) Execute(req UpdateUserProfileRequest) (*UpdateUserProfileResponse, error) {
	// 1. Validate user ID
	if req.UserID == "" {
		return nil, ErrUserNotFoundUpdate
	}

	// 2. Get existing user
	user, err := uc.UserRepository.GetUserByID(req.UserID)
	if err != nil {
		return nil, ErrUserNotFoundUpdate
	}

	// 3. Parse and validate BirthDate
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, ErrInvalidBirthDateUpdate
	}

	// 4. Update user fields
	user.Name = req.Name
	user.BirthDate = birthDate

	// 5. Validate user fields using domain methods
	if !user.IsValidName() {
		return nil, ErrInvalidNameUpdate
	}

	if !user.IsPastDate() {
		return nil, ErrFutureBirthDateUpdate
	}

	// 6. Save updated user
	err = uc.UserRepository.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	// 7. Return response
	return &UpdateUserProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		BirthDate: user.BirthDate.Format("2006-01-02"),
	}, nil
}

package usecase

import (
	"errors"
	"time"

	"github.com/example/cadastro-de-usuarios/domain"
)

var (
	ErrInvalidNameUpdate      = errors.New("nome deve ter entre 2 e 50 caracteres e conter apenas letras e espaços")
	ErrInvalidBirthDateUpdate = errors.New("data de nascimento inválida")
	ErrFutureBirthDateUpdate  = errors.New("data de nascimento não pode ser no futuro")
	ErrUserNotFoundUpdate     = errors.New("usuário não encontrado")
)

type UpdateUserProfileInput struct {
	UserID    string `param:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birthDate"`
}

type UpdateUserProfileOutput struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	BirthDate string `json:"birthDate"`
}

type UpdateUserProfileUseCase struct {
	UserRepository domain.UserRepository
}

func NewUpdateUserProfileUseCase(repo domain.UserRepository) *UpdateUserProfileUseCase {
	return &UpdateUserProfileUseCase{
		UserRepository: repo,
	}
}

func (uc *UpdateUserProfileUseCase) validateUserID(userID string) error {
	if userID == "" {
		return ErrUserNotFoundUpdate
	}
	return nil
}

func (uc *UpdateUserProfileUseCase) parseBirthDate(birthDateStr string) (time.Time, error) {
	birthDate, parseErr := time.Parse("2006-01-02", birthDateStr)
	if parseErr != nil {
		return time.Time{}, ErrInvalidBirthDateUpdate
	}
	return birthDate, nil
}

func (uc *UpdateUserProfileUseCase) updateUserFields(user *domain.User, name string, birthDate time.Time) error {
	user.Name = name
	user.BirthDate = birthDate

	if !user.IsValidName() {
		return ErrInvalidNameUpdate
	}

	if !user.IsPastDate() {
		return ErrFutureBirthDateUpdate
	}

	return nil
}

func (uc *UpdateUserProfileUseCase) saveUser(user *domain.User) error {
	updateErr := uc.UserRepository.UpdateUser(user)
	if updateErr != nil {
		return updateErr
	}
	return nil
}

func (uc *UpdateUserProfileUseCase) Execute(req UpdateUserProfileInput) (*UpdateUserProfileOutput, error) {

	if userIDValidationErr := uc.validateUserID(req.UserID); userIDValidationErr != nil {
		return nil, userIDValidationErr
	}

	user, repositoryErr := uc.UserRepository.GetUserByID(req.UserID)
	if repositoryErr != nil {
		return nil, ErrUserNotFoundUpdate
	}

	birthDate, parseErr := uc.parseBirthDate(req.BirthDate)
	if parseErr != nil {
		return nil, parseErr
	}

	if fieldUpdateErr := uc.updateUserFields(user, req.Name, birthDate); fieldUpdateErr != nil {
		return nil, fieldUpdateErr
	}

	if saveErr := uc.saveUser(user); saveErr != nil {
		return nil, saveErr
	}

	return &UpdateUserProfileOutput{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		BirthDate: user.BirthDate.Format("2006-01-02"),
	}, nil
}

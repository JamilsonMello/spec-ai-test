package usecase

import (
	"errors"
	"time"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrInvalidNameUpdate      = errors.New("nome deve ter entre 2 e 50 caracteres e conter apenas letras e espaços")
	ErrInvalidBirthDateUpdate = errors.New("data de nascimento inválida")
	ErrFutureBirthDateUpdate  = errors.New("data de nascimento não pode ser no futuro")
	ErrUserNotFoundUpdate     = errors.New("usuário não encontrado")
)

type UpdateUserProfileRequest struct {
	UserID    string `param:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birthDate"`
}

type UpdateUserProfileResponse struct {
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

func (uc *UpdateUserProfileUseCase) Execute(req UpdateUserProfileRequest) (*UpdateUserProfileResponse, error) {
	user, err := uc.getUser(req.UserID)
	if err != nil {
		return nil, err
	}

	birthDate, err := uc.parseBirthDate(req.BirthDate)
	if err != nil {
		return nil, err
	}

	user.Name = req.Name
	user.BirthDate = birthDate

	if err := uc.validateUpdatedUser(user); err != nil {
		return nil, err
	}

	err = uc.UserRepository.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return uc.buildResponse(user), nil
}

func (uc *UpdateUserProfileUseCase) getUser(userID string) (*domain.User, error) {
	if userID == "" {
		return nil, ErrUserNotFoundUpdate
	}

	user, err := uc.UserRepository.FindUserByUuid(userID)
	if err != nil {
		return nil, ErrUserNotFoundUpdate
	}
	return user, nil
}

func (uc *UpdateUserProfileUseCase) parseBirthDate(birthDateStr string) (time.Time, error) {
	birthDate, err := time.Parse("2006-01-02", birthDateStr)
	if err != nil {
		return time.Time{}, ErrInvalidBirthDateUpdate
	}
	return birthDate, nil
}

func (uc *UpdateUserProfileUseCase) validateUpdatedUser(user *domain.User) error {
	if !user.IsValidName() {
		return ErrInvalidNameUpdate
	}

	if !user.IsPastDate() {
		return ErrFutureBirthDateUpdate
	}

	return nil
}

func (uc *UpdateUserProfileUseCase) buildResponse(user *domain.User) *UpdateUserProfileResponse {
	return &UpdateUserProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		BirthDate: user.BirthDate.Format("2006-01-02"),
	}
}

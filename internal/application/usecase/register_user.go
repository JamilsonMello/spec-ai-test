package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrInvalidName      = errors.New("nome deve ter entre 2 e 50 caracteres e conter apenas letras e espaços")
	ErrInvalidSurname   = errors.New("sobrenome deve ter entre 2 e 50 caracteres e conter apenas letras e espaços")
	ErrInvalidEmail     = errors.New("email inválido")
	ErrEmailInUse       = errors.New("email já está em uso")
	ErrInvalidBirthDate = errors.New("data de nascimento inválida")
	ErrUserTooYoung     = errors.New("usuário deve ter no mínimo 18 anos")
	ErrFutureBirthDate  = errors.New("data de nascimento não pode ser no futuro")
)

type RegisterUserRequest struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
	BirthDate string `json:"birthDate"`
}

type RegisterUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	BirthDate string    `json:"birthDate"`
}

type RegisterUserUseCase struct {
	UserRepository domain.UserRepository
}

func NewRegisterUserUseCase(repo domain.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		UserRepository: repo,
	}
}

func (uc *RegisterUserUseCase) Execute(req RegisterUserRequest) (*RegisterUserResponse, error) {
	user, err := uc.createUserFromRequest(req)
	if err != nil {
		return nil, err
	}

	if err := uc.validateUser(user); err != nil {
		return nil, err
	}

	if err := uc.checkEmailUniqueness(user.Email); err != nil {
		return nil, err
	}

	user.ID = uuid.New().String()

	err = uc.UserRepository.SaveUser(user)
	if err != nil {
		return nil, err
	}

	return uc.buildResponse(user), nil
}

func (uc *RegisterUserUseCase) createUserFromRequest(req RegisterUserRequest) (*domain.User, error) {
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, ErrInvalidBirthDate
	}

	return &domain.User{
		Name:      req.Name,
		Surname:   req.Surname,
		Email:     req.Email,
		BirthDate: birthDate,
		Role:      "user",
		CreatedAt: time.Now(),
	}, nil
}

func (uc *RegisterUserUseCase) validateUser(user *domain.User) error {
	if !user.IsValidName() {
		return ErrInvalidName
	}

	if !user.IsValidSurname() {
		return ErrInvalidSurname
	}

	if !user.IsValidEmailFormat() {
		return ErrInvalidEmail
	}

	if !user.IsPastDate() {
		return ErrFutureBirthDate
	}

	if !user.IsAdult() {
		return ErrUserTooYoung
	}

	return nil
}

func (uc *RegisterUserUseCase) checkEmailUniqueness(email string) error {
	existingUser, err := uc.UserRepository.GetUserByEmail(email)
	if err != nil && err != errors.New("user not found") {
		return err
	}
	if existingUser != nil {
		return ErrEmailInUse
	}
	return nil
}

func (uc *RegisterUserUseCase) buildResponse(user *domain.User) *RegisterUserResponse {
	return &RegisterUserResponse{
		ID:        uuid.MustParse(user.ID),
		Name:      user.Name,
		Surname:   user.Surname,
		Email:     user.Email,
		BirthDate: user.BirthDate.Format("2006-01-02"),
	}
}

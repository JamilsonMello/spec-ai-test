package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/domain"
)

// Custom errors for use case validation
var (
	ErrInvalidName       = errors.New("nome deve ter entre 2 e 50 caracteres e conter apenas letras e espaços")
	ErrInvalidSurname    = errors.New("sobrenome deve ter entre 2 e 50 caracteres e conter apenas letras e espaços")
	ErrInvalidEmail      = errors.New("email inválido")
	ErrEmailInUse        = errors.New("email já está em uso")
	ErrInvalidBirthDate  = errors.New("data de nascimento inválida")
	ErrUserTooYoung      = errors.New("usuário deve ter no mínimo 18 anos")
	ErrFutureBirthDate   = errors.New("data de nascimento não pode ser no futuro")
)

// RegisterUserUseCase handles the business logic for user registration.
type RegisterUserUseCase struct {
	UserRepository domain.UserRepository
}

// NewRegisterUserUseCase creates a new RegisterUserUseCase.
func NewRegisterUserUseCase(repo domain.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		UserRepository: repo,
	}
}

// Execute performs the user registration process.
func (uc *RegisterUserUseCase) Execute(req RegisterUserRequest) (*RegisterUserResponse, error) {
	// 1. Parse and validate BirthDate
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, ErrInvalidBirthDate
	}

	user := &domain.User{
		Name:      req.Name,
		Surname:   req.Surname,
		Email:     req.Email,
		BirthDate: birthDate,
		Role:      "user",
		CreatedAt: time.Now(),
	}

	// 2. Validate user fields using domain methods
	if !user.IsValidName() {
		return nil, ErrInvalidName
	}

	if !user.IsValidSurname() {
		return nil, ErrInvalidSurname
	}

	if !user.IsValidEmailFormat() {
		return nil, ErrInvalidEmail
	}

	if !user.IsPastDate() {
		return nil, ErrFutureBirthDate
	}

	if !user.IsAdult() {
		return nil, ErrUserTooYoung
	}

	// 3. Check for email uniqueness
	existingUser, err := uc.UserRepository.GetUserByEmail(user.Email)
	if err != nil && err != errors.New("user not found") { // Handle actual errors other than not found
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailInUse
	}

	// 4. Generate ID and save user
	user.ID = uuid.New().String()

	err = uc.UserRepository.SaveUser(user)
	if err != nil {
		return nil, err
	}

	// 5. Return response
	return &RegisterUserResponse{
		ID:          uuid.MustParse(user.ID),
		Name:        user.Name,
		Surname:     user.Surname,
		Email:       user.Email,
		BirthDate:   user.BirthDate.Format("2006-01-02"),
	}, nil
}

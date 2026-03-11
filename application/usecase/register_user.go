package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/domain"
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

type RegisterUserInput struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
	BirthDate string `json:"birthDate"`
}

type RegisterUserOutput struct {
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

func (uc *RegisterUserUseCase) validateUserInput(user *domain.User) error {
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
	existingUser, lookupErr := uc.UserRepository.GetUserByEmail(email)
	if lookupErr != nil && lookupErr != errors.New("user not found") {
		return lookupErr
	}
	if existingUser != nil {
		return ErrEmailInUse
	}
	return nil
}

func (uc *RegisterUserUseCase) parseBirthDate(birthDateStr string) (time.Time, error) {
	birthDate, parseErr := time.Parse("2006-01-02", birthDateStr)
	if parseErr != nil {
		return time.Time{}, ErrInvalidBirthDate
	}
	return birthDate, nil
}

func (uc *RegisterUserUseCase) createUser(req RegisterUserInput, birthDate time.Time) *domain.User {
	return &domain.User{
		Name:      req.Name,
		Surname:   req.Surname,
		Email:     req.Email,
		BirthDate: birthDate,
		Role:      "user",
		CreatedAt: time.Now(),
	}
}

func (uc *RegisterUserUseCase) assignUserID(user *domain.User) {
	user.ID = uuid.New().String()
}

func (uc *RegisterUserUseCase) saveUserToRepo(user *domain.User) error {
	saveErr := uc.UserRepository.SaveUser(user)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

func (uc *RegisterUserUseCase) mapUserToOutput(user *domain.User) *RegisterUserOutput {
	return &RegisterUserOutput{
		ID:        uuid.MustParse(user.ID),
		Name:      user.Name,
		Surname:   user.Surname,
		Email:     user.Email,
		BirthDate: user.BirthDate.Format("2006-01-02"),
	}
}

func (uc *RegisterUserUseCase) Execute(req RegisterUserInput) (*RegisterUserOutput, error) {
	birthDate, birthDateParseErr := uc.parseBirthDate(req.BirthDate)
	if birthDateParseErr != nil {
		return nil, birthDateParseErr
	}

	user := uc.createUser(req, birthDate)

	if validationErr := uc.validateUserInput(user); validationErr != nil {
		return nil, validationErr
	}

	if emailCheckErr := uc.checkEmailUniqueness(user.Email); emailCheckErr != nil {
		return nil, emailCheckErr
	}

	uc.assignUserID(user)

	if saveErr := uc.saveUserToRepo(user); saveErr != nil {
		return nil, saveErr
	}

	return uc.mapUserToOutput(user), nil
}

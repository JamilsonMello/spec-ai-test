package usecase

import (
	"errors"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrInvalidCredentials = errors.New("Credenciais inválidas")
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type TokenGenerator interface {
	GenerateToken(userID string) (string, error)
}

type LoginUseCase struct {
	UserRepository domain.UserRepository
	PasswordHasher domain.PasswordHasher
	TokenGenerator TokenGenerator
}

func NewLoginUseCase(repo domain.UserRepository, hasher domain.PasswordHasher, tokenGen TokenGenerator) *LoginUseCase {
	return &LoginUseCase{
		UserRepository: repo,
		PasswordHasher: hasher,
		TokenGenerator: tokenGen,
	}
}

func (uc *LoginUseCase) Execute(req LoginRequest) (*LoginResponse, error) {
	user, err := uc.UserRepository.GetUserByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := uc.PasswordHasher.Compare(user.Password, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := uc.TokenGenerator.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Token: token}, nil
}

package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) GenerateToken(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func newLoginTestDeps() (*MockUserRepository, *MockPasswordHasher, *MockTokenGenerator, *LoginUseCase) {
	mockUserRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	mockTokenGen := new(MockTokenGenerator)

	uc := NewLoginUseCase(mockUserRepo, mockHasher, mockTokenGen)
	return mockUserRepo, mockHasher, mockTokenGen, uc
}

func TestLogin_Success(t *testing.T) {
	mockUserRepo, mockHasher, mockTokenGen, uc := newLoginTestDeps()

	user := &domain.User{ID: "user-123", Email: "test@example.com", Password: "$2a$10$hashed"}

	mockUserRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockHasher.On("Compare", "$2a$10$hashed", "password123").Return(nil)
	mockTokenGen.On("GenerateToken", "user-123").Return("jwt-token", nil)

	resp, err := uc.Execute(LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "jwt-token", resp.Token)
	mockUserRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockUserRepo, _, _, uc := newLoginTestDeps()

	mockUserRepo.On("GetUserByEmail", "notfound@example.com").Return(nil, domain.ErrUserNotFound)

	resp, err := uc.Execute(LoginRequest{
		Email:    "notfound@example.com",
		Password: "password123",
	})

	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockUserRepo, mockHasher, _, uc := newLoginTestDeps()

	user := &domain.User{ID: "user-123", Email: "test@example.com", Password: "$2a$10$hashed"}

	mockUserRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockHasher.On("Compare", "$2a$10$hashed", "wrongpassword").Return(errors.New("mismatch"))

	resp, err := uc.Execute(LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})

	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestLogin_TokenGenerationError(t *testing.T) {
	mockUserRepo, mockHasher, mockTokenGen, uc := newLoginTestDeps()

	user := &domain.User{ID: "user-123", Email: "test@example.com", Password: "$2a$10$hashed"}

	mockUserRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockHasher.On("Compare", "$2a$10$hashed", "password123").Return(nil)
	mockTokenGen.On("GenerateToken", "user-123").Return("", errors.New("token error"))

	resp, err := uc.Execute(LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.Nil(t, resp)
	assert.EqualError(t, err, "token error")
}

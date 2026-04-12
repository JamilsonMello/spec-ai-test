package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type MockPasswordRecoveryRepository struct {
	mock.Mock
}

func (m *MockPasswordRecoveryRepository) SavePasswordRecovery(recovery *domain.PasswordRecovery) error {
	args := m.Called(recovery)
	return args.Error(0)
}

func (m *MockPasswordRecoveryRepository) GetPasswordRecoveryByToken(token string) (*domain.PasswordRecovery, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PasswordRecovery), args.Error(1)
}

func (m *MockPasswordRecoveryRepository) UpdatePasswordRecovery(recovery *domain.PasswordRecovery) error {
	args := m.Called(recovery)
	return args.Error(0)
}

func (m *MockPasswordRecoveryRepository) UpdatePasswordRecoveryTx(tx interface{}, recovery *domain.PasswordRecovery) error {
	args := m.Called(tx, recovery)
	return args.Error(0)
}

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendPasswordRecoveryEmail(email string, token string) error {
	args := m.Called(email, token)
	return args.Error(0)
}

func TestRequestPasswordRecovery_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRecoveryRepo := new(MockPasswordRecoveryRepository)
	mockEmailSender := new(MockEmailSender)

	uc := NewRequestPasswordRecoveryUseCase(mockUserRepo, mockRecoveryRepo, mockEmailSender)

	user := &domain.User{ID: "user-123", Email: "test@example.com"}
	mockUserRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockRecoveryRepo.On("SavePasswordRecovery", mock.AnythingOfType("*domain.PasswordRecovery")).Return(nil)
	mockEmailSender.On("SendPasswordRecoveryEmail", "test@example.com", mock.AnythingOfType("string")).Return(nil)

	resp, err := uc.Execute(RequestPasswordRecoveryRequest{Email: "test@example.com"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Se um usuário com esse e-mail existir, você receberá instruções em breve.", resp.Message)
	mockUserRepo.AssertExpectations(t)
	mockRecoveryRepo.AssertExpectations(t)
	mockEmailSender.AssertExpectations(t)
}

func TestRequestPasswordRecovery_UserNotFound_ReturnsGenericMessage(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRecoveryRepo := new(MockPasswordRecoveryRepository)
	mockEmailSender := new(MockEmailSender)

	uc := NewRequestPasswordRecoveryUseCase(mockUserRepo, mockRecoveryRepo, mockEmailSender)

	mockUserRepo.On("GetUserByEmail", "nonexistent@example.com").Return(nil, domain.ErrUserNotFound)

	resp, err := uc.Execute(RequestPasswordRecoveryRequest{Email: "nonexistent@example.com"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Se um usuário com esse e-mail existir, você receberá instruções em breve.", resp.Message)
	mockRecoveryRepo.AssertNotCalled(t, "SavePasswordRecovery")
	mockEmailSender.AssertNotCalled(t, "SendPasswordRecoveryEmail")
}

func TestRequestPasswordRecovery_SaveRecoveryError_ReturnsGenericMessage(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRecoveryRepo := new(MockPasswordRecoveryRepository)
	mockEmailSender := new(MockEmailSender)

	uc := NewRequestPasswordRecoveryUseCase(mockUserRepo, mockRecoveryRepo, mockEmailSender)

	user := &domain.User{ID: "user-123", Email: "test@example.com"}
	mockUserRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockRecoveryRepo.On("SavePasswordRecovery", mock.AnythingOfType("*domain.PasswordRecovery")).Return(errors.New("db error"))

	resp, err := uc.Execute(RequestPasswordRecoveryRequest{Email: "test@example.com"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Se um usuário com esse e-mail existir, você receberá instruções em breve.", resp.Message)
	mockEmailSender.AssertNotCalled(t, "SendPasswordRecoveryEmail")
}

func TestRequestPasswordRecovery_EmailSendError_StillReturnsSuccess(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRecoveryRepo := new(MockPasswordRecoveryRepository)
	mockEmailSender := new(MockEmailSender)

	uc := NewRequestPasswordRecoveryUseCase(mockUserRepo, mockRecoveryRepo, mockEmailSender)

	user := &domain.User{ID: "user-123", Email: "test@example.com"}
	mockUserRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockRecoveryRepo.On("SavePasswordRecovery", mock.AnythingOfType("*domain.PasswordRecovery")).Return(nil)
	mockEmailSender.On("SendPasswordRecoveryEmail", "test@example.com", mock.AnythingOfType("string")).Return(errors.New("smtp error"))

	resp, err := uc.Execute(RequestPasswordRecoveryRequest{Email: "test@example.com"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Se um usuário com esse e-mail existir, você receberá instruções em breve.", resp.Message)
}

func TestRequestPasswordRecovery_NilEmailSender_StillReturnsSuccess(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRecoveryRepo := new(MockPasswordRecoveryRepository)

	uc := NewRequestPasswordRecoveryUseCase(mockUserRepo, mockRecoveryRepo, nil)

	user := &domain.User{ID: "user-123", Email: "test@example.com"}
	mockUserRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockRecoveryRepo.On("SavePasswordRecovery", mock.AnythingOfType("*domain.PasswordRecovery")).Return(nil)

	resp, err := uc.Execute(RequestPasswordRecoveryRequest{Email: "test@example.com"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Se um usuário com esse e-mail existir, você receberá instruções em breve.", resp.Message)
}

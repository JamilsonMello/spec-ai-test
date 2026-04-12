package usecase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

func (m *MockUserRepository) UpdateUserTx(tx interface{}, user *domain.User) error {
	args := m.Called(tx, user)
	return args.Error(0)
}

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) Compare(hashedPassword string, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) ExecuteInTransaction(fn func(tx interface{}) error) error {
	args := m.Called(mock.AnythingOfType("func(interface {}) error"))
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return fn("mock-tx")
}

func newResetPasswordTestDeps() (*MockUserRepository, *MockPasswordRecoveryRepository, *MockPasswordHasher, *MockTransactionManager, *ResetPasswordUseCase) {
	mockUserRepo := new(MockUserRepository)
	mockRecoveryRepo := new(MockPasswordRecoveryRepository)
	mockHasher := new(MockPasswordHasher)
	mockTxManager := new(MockTransactionManager)

	uc := NewResetPasswordUseCase(mockUserRepo, mockUserRepo, mockRecoveryRepo, mockRecoveryRepo, mockHasher, mockTxManager)
	return mockUserRepo, mockRecoveryRepo, mockHasher, mockTxManager, uc
}

func TestResetPassword_Success(t *testing.T) {
	mockUserRepo, mockRecoveryRepo, mockHasher, mockTxManager, uc := newResetPasswordTestDeps()

	recovery := &domain.PasswordRecovery{
		ID:        "rec-123",
		Token:     "valid-token",
		UserID:    "user-123",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}
	user := &domain.User{ID: "user-123", Email: "test@example.com"}

	mockRecoveryRepo.On("GetPasswordRecoveryByToken", "valid-token").Return(recovery, nil)
	mockUserRepo.On("FindUserByUuid", "user-123").Return(user, nil)
	mockHasher.On("Hash", "newpassword123").Return("hashed-password", nil)
	mockTxManager.On("ExecuteInTransaction", mock.Anything).Return(nil)
	mockUserRepo.On("UpdateUserTx", "mock-tx", mock.AnythingOfType("*domain.User")).Return(nil)
	mockRecoveryRepo.On("UpdatePasswordRecoveryTx", "mock-tx", mock.AnythingOfType("*domain.PasswordRecovery")).Return(nil)

	resp, err := uc.Execute(ResetPasswordRequest{
		Token:       "valid-token",
		NewPassword: "newpassword123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Senha alterada com sucesso.", resp.Message)
	mockRecoveryRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestResetPassword_TokenExpired(t *testing.T) {
	_, mockRecoveryRepo, _, _, uc := newResetPasswordTestDeps()

	recovery := &domain.PasswordRecovery{
		ID:        "rec-123",
		Token:     "expired-token",
		UserID:    "user-123",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		Used:      false,
	}

	mockRecoveryRepo.On("GetPasswordRecoveryByToken", "expired-token").Return(recovery, nil)

	resp, err := uc.Execute(ResetPasswordRequest{
		Token:       "expired-token",
		NewPassword: "newpassword123",
	})

	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestResetPassword_TokenAlreadyUsed(t *testing.T) {
	_, mockRecoveryRepo, _, _, uc := newResetPasswordTestDeps()

	recovery := &domain.PasswordRecovery{
		ID:        "rec-123",
		Token:     "used-token",
		UserID:    "user-123",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      true,
	}

	mockRecoveryRepo.On("GetPasswordRecoveryByToken", "used-token").Return(recovery, nil)

	resp, err := uc.Execute(ResetPasswordRequest{
		Token:       "used-token",
		NewPassword: "newpassword123",
	})

	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestResetPassword_TokenNotFound(t *testing.T) {
	_, mockRecoveryRepo, _, _, uc := newResetPasswordTestDeps()

	mockRecoveryRepo.On("GetPasswordRecoveryByToken", "invalid-token").Return(nil, domain.ErrRecoveryTokenNotFound)

	resp, err := uc.Execute(ResetPasswordRequest{
		Token:       "invalid-token",
		NewPassword: "newpassword123",
	})

	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestResetPassword_PasswordTooShort(t *testing.T) {
	_, _, _, _, uc := newResetPasswordTestDeps()

	resp, err := uc.Execute(ResetPasswordRequest{
		Token:       "some-token",
		NewPassword: "short",
	})

	assert.Nil(t, resp)
	assert.Equal(t, ErrPasswordTooShort, err)
}

func TestResetPassword_UserNotFound(t *testing.T) {
	mockUserRepo, mockRecoveryRepo, _, _, uc := newResetPasswordTestDeps()

	recovery := &domain.PasswordRecovery{
		ID:        "rec-123",
		Token:     "valid-token",
		UserID:    "nonexistent-user",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	mockRecoveryRepo.On("GetPasswordRecoveryByToken", "valid-token").Return(recovery, nil)
	mockUserRepo.On("FindUserByUuid", "nonexistent-user").Return(nil, domain.ErrUserNotFound)

	resp, err := uc.Execute(ResetPasswordRequest{
		Token:       "valid-token",
		NewPassword: "newpassword123",
	})

	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestResetPassword_UsesPasswordHasher(t *testing.T) {
	mockUserRepo, mockRecoveryRepo, mockHasher, mockTxManager, uc := newResetPasswordTestDeps()

	recovery := &domain.PasswordRecovery{
		ID:        "rec-123",
		Token:     "valid-token",
		UserID:    "user-123",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}
	user := &domain.User{ID: "user-123", Email: "test@example.com"}

	mockRecoveryRepo.On("GetPasswordRecoveryByToken", "valid-token").Return(recovery, nil)
	mockUserRepo.On("FindUserByUuid", "user-123").Return(user, nil)
	mockHasher.On("Hash", "securepassword").Return("$2a$10$hashed", nil)
	mockTxManager.On("ExecuteInTransaction", mock.Anything).Return(nil)
	mockUserRepo.On("UpdateUserTx", "mock-tx", mock.MatchedBy(func(u *domain.User) bool {
		return u.Password == "$2a$10$hashed"
	})).Return(nil)
	mockRecoveryRepo.On("UpdatePasswordRecoveryTx", "mock-tx", mock.MatchedBy(func(r *domain.PasswordRecovery) bool {
		return r.Used == true
	})).Return(nil)

	resp, err := uc.Execute(ResetPasswordRequest{
		Token:       "valid-token",
		NewPassword: "securepassword",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockHasher.AssertCalled(t, "Hash", "securepassword")
}

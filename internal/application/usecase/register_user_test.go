package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type mockUserRepository struct {
	saveUserFn              func(user *domain.User) error
	getUserByEmailFn        func(email string) (*domain.User, error)
	findUserByUuidFn        func(id string) (*domain.User, error)
	deleteUserFn            func(id string) error
	updateUserFn            func(user *domain.User) error
	updateProfilePictureURL func(userID string, url string) error
	listUsersFn             func(filter domain.UserFilter, page int, limit int) ([]*domain.User, int, error)
}

func (m *mockUserRepository) SaveUser(user *domain.User) error {
	return m.saveUserFn(user)
}

func (m *mockUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	return m.getUserByEmailFn(email)
}

func (m *mockUserRepository) FindUserByUuid(id string) (*domain.User, error) {
	return m.findUserByUuidFn(id)
}

func (m *mockUserRepository) DeleteUser(id string) error {
	return m.deleteUserFn(id)
}

func (m *mockUserRepository) UpdateUser(user *domain.User) error {
	return m.updateUserFn(user)
}

func (m *mockUserRepository) UpdateProfilePictureURL(userID string, url string) error {
	return m.updateProfilePictureURL(userID, url)
}

func (m *mockUserRepository) ListUsers(filter domain.UserFilter, page int, limit int) ([]*domain.User, int, error) {
	return m.listUsersFn(filter, page, limit)
}

func newSuccessMock() *mockUserRepository {
	return &mockUserRepository{
		saveUserFn: func(user *domain.User) error {
			return nil
		},
		getUserByEmailFn: func(email string) (*domain.User, error) {
			return nil, nil
		},
	}
}

func validRequest() RegisterUserRequest {
	birthDate := time.Now().AddDate(-20, 0, 0).Format("2006-01-02")
	return RegisterUserRequest{
		Name:      "Joao",
		Surname:   "Silva",
		Email:     "joao@email.com",
		BirthDate: birthDate,
	}
}

func TestRegisterUserUseCase_Execute(t *testing.T) {
	t.Run("sucesso ao registrar usuario valido", func(t *testing.T) {
		mock := newSuccessMock()
		uc := NewRegisterUserUseCase(mock)

		resp, err := uc.Execute(validRequest())

		if err != nil {
			t.Fatalf("esperava sucesso, obteve erro: %v", err)
		}
		if resp == nil {
			t.Fatal("esperava resposta, obteve nil")
		}
		if resp.Name != "Joao" {
			t.Errorf("esperava Name 'Joao', obteve '%s'", resp.Name)
		}
		if resp.Email != "joao@email.com" {
			t.Errorf("esperava Email 'joao@email.com', obteve '%s'", resp.Email)
		}
	})

	t.Run("erro ao nome ser vazio", func(t *testing.T) {
		mock := newSuccessMock()
		uc := NewRegisterUserUseCase(mock)

		req := validRequest()
		req.Name = ""

		_, err := uc.Execute(req)

		if !errors.Is(err, ErrInvalidName) {
			t.Errorf("esperava ErrInvalidName, obteve: %v", err)
		}
	})

	t.Run("erro ao nome ter menos de 2 caracteres", func(t *testing.T) {
		mock := newSuccessMock()
		uc := NewRegisterUserUseCase(mock)

		req := validRequest()
		req.Name = "A"

		_, err := uc.Execute(req)

		if !errors.Is(err, ErrInvalidName) {
			t.Errorf("esperava ErrInvalidName, obteve: %v", err)
		}
	})

	t.Run("erro ao sobrenome ser vazio", func(t *testing.T) {
		mock := newSuccessMock()
		uc := NewRegisterUserUseCase(mock)

		req := validRequest()
		req.Surname = ""

		_, err := uc.Execute(req)

		if !errors.Is(err, ErrInvalidSurname) {
			t.Errorf("esperava ErrInvalidSurname, obteve: %v", err)
		}
	})

	t.Run("erro ao email ser invalido", func(t *testing.T) {
		mock := newSuccessMock()
		uc := NewRegisterUserUseCase(mock)

		req := validRequest()
		req.Email = "email-invalido"

		_, err := uc.Execute(req)

		if !errors.Is(err, ErrInvalidEmail) {
			t.Errorf("esperava ErrInvalidEmail, obteve: %v", err)
		}
	})

	t.Run("erro ao data de nascimento ser no futuro", func(t *testing.T) {
		mock := newSuccessMock()
		uc := NewRegisterUserUseCase(mock)

		req := validRequest()
		req.BirthDate = time.Now().AddDate(1, 0, 0).Format("2006-01-02")

		_, err := uc.Execute(req)

		if !errors.Is(err, ErrFutureBirthDate) {
			t.Errorf("esperava ErrFutureBirthDate, obteve: %v", err)
		}
	})

	t.Run("erro ao usuario ser menor de 18 anos", func(t *testing.T) {
		mock := newSuccessMock()
		uc := NewRegisterUserUseCase(mock)

		req := validRequest()
		req.BirthDate = time.Now().AddDate(-17, 0, 0).Format("2006-01-02")

		_, err := uc.Execute(req)

		if !errors.Is(err, ErrUserTooYoung) {
			t.Errorf("esperava ErrUserTooYoung, obteve: %v", err)
		}
	})

	t.Run("erro ao data de nascimento ter formato invalido", func(t *testing.T) {
		mock := newSuccessMock()
		uc := NewRegisterUserUseCase(mock)

		req := validRequest()
		req.BirthDate = "invalido"

		_, err := uc.Execute(req)

		if !errors.Is(err, ErrInvalidBirthDate) {
			t.Errorf("esperava ErrInvalidBirthDate, obteve: %v", err)
		}
	})

	t.Run("erro ao email ja estar em uso", func(t *testing.T) {
		mock := newSuccessMock()
		mock.getUserByEmailFn = func(email string) (*domain.User, error) {
			return &domain.User{ID: "existing-id", Email: email}, nil
		}
		uc := NewRegisterUserUseCase(mock)

		_, err := uc.Execute(validRequest())

		if !errors.Is(err, ErrEmailInUse) {
			t.Errorf("esperava ErrEmailInUse, obteve: %v", err)
		}
	})

	t.Run("erro ao repositorio falhar no SaveUser", func(t *testing.T) {
		repoErr := errors.New("database connection failed")
		mock := newSuccessMock()
		mock.saveUserFn = func(user *domain.User) error {
			return repoErr
		}
		uc := NewRegisterUserUseCase(mock)

		_, err := uc.Execute(validRequest())

		if err == nil {
			t.Fatal("esperava erro, obteve nil")
		}
		if err.Error() != repoErr.Error() {
			t.Errorf("esperava erro '%s', obteve: '%s'", repoErr.Error(), err.Error())
		}
	})

	t.Run("erro ao repositorio falhar no GetUserByEmail", func(t *testing.T) {
		repoErr := errors.New("database connection failed")
		mock := newSuccessMock()
		mock.getUserByEmailFn = func(email string) (*domain.User, error) {
			return nil, repoErr
		}
		uc := NewRegisterUserUseCase(mock)

		_, err := uc.Execute(validRequest())

		if err == nil {
			t.Fatal("esperava erro, obteve nil")
		}
	})
}

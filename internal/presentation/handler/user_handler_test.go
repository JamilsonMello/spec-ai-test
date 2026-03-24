package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type mockUserRepo struct {
	saveUserFn              func(user *domain.User) error
	getUserByEmailFn        func(email string) (*domain.User, error)
	findUserByUuidFn        func(id string) (*domain.User, error)
	deleteUserFn            func(id string) error
	updateUserFn            func(user *domain.User) error
	updateProfilePictureURL func(userID string, url string) error
	listUsersFn             func(filter domain.UserFilter, page int, limit int) ([]*domain.User, int, error)
}

func (m *mockUserRepo) SaveUser(user *domain.User) error {
	return m.saveUserFn(user)
}

func (m *mockUserRepo) GetUserByEmail(email string) (*domain.User, error) {
	return m.getUserByEmailFn(email)
}

func (m *mockUserRepo) FindUserByUuid(id string) (*domain.User, error) {
	return m.findUserByUuidFn(id)
}

func (m *mockUserRepo) DeleteUser(id string) error {
	return m.deleteUserFn(id)
}

func (m *mockUserRepo) UpdateUser(user *domain.User) error {
	return m.updateUserFn(user)
}

func (m *mockUserRepo) UpdateProfilePictureURL(userID string, url string) error {
	return m.updateProfilePictureURL(userID, url)
}

func (m *mockUserRepo) ListUsers(filter domain.UserFilter, page int, limit int) ([]*domain.User, int, error) {
	return m.listUsersFn(filter, page, limit)
}

func newHandlerSuccessMock() *mockUserRepo {
	return &mockUserRepo{
		saveUserFn: func(user *domain.User) error {
			return nil
		},
		getUserByEmailFn: func(email string) (*domain.User, error) {
			return nil, nil
		},
	}
}

func newTestHandler(repo domain.UserRepository) *UserHandler {
	registerUC := usecase.NewRegisterUserUseCase(repo)
	return &UserHandler{
		RegisterUserUseCase: registerUC,
	}
}

func validJSON() string {
	birthDate := time.Now().AddDate(-20, 0, 0).Format("2006-01-02")
	return `{"name":"Joao","surname":"Silva","email":"joao@email.com","birthDate":"` + birthDate + `"}`
}

func executeRegisterRequest(handler *UserHandler, body string) *httptest.ResponseRecorder {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/usuarios", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler.RegisterUser(c)
	return rec
}

func TestRegisterUserHandler(t *testing.T) {
	t.Run("sucesso ao registrar usuario - status 200", func(t *testing.T) {
		mock := newHandlerSuccessMock()
		handler := newTestHandler(mock)

		rec := executeRegisterRequest(handler, validJSON())

		if rec.Code != http.StatusOK {
			t.Errorf("esperava status %d, obteve %d", http.StatusOK, rec.Code)
		}

		var resp usecase.RegisterUserResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("erro ao decodificar resposta: %v", err)
		}
		if resp.Name != "Joao" {
			t.Errorf("esperava Name 'Joao', obteve '%s'", resp.Name)
		}
	})

	t.Run("erro 400 ao enviar JSON invalido", func(t *testing.T) {
		mock := newHandlerSuccessMock()
		handler := newTestHandler(mock)

		rec := executeRegisterRequest(handler, "{invalid json")

		if rec.Code != http.StatusBadRequest {
			t.Errorf("esperava status %d, obteve %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("erro 422 ao nome ser vazio", func(t *testing.T) {
		mock := newHandlerSuccessMock()
		handler := newTestHandler(mock)

		birthDate := time.Now().AddDate(-20, 0, 0).Format("2006-01-02")
		body := `{"name":"","surname":"Silva","email":"joao@email.com","birthDate":"` + birthDate + `"}`

		rec := executeRegisterRequest(handler, body)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("esperava status %d, obteve %d", http.StatusUnprocessableEntity, rec.Code)
		}

		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		if resp["error"] != usecase.ErrInvalidName.Error() {
			t.Errorf("esperava erro '%s', obteve '%s'", usecase.ErrInvalidName.Error(), resp["error"])
		}
	})

	t.Run("erro 422 ao sobrenome ser vazio", func(t *testing.T) {
		mock := newHandlerSuccessMock()
		handler := newTestHandler(mock)

		birthDate := time.Now().AddDate(-20, 0, 0).Format("2006-01-02")
		body := `{"name":"Joao","surname":"","email":"joao@email.com","birthDate":"` + birthDate + `"}`

		rec := executeRegisterRequest(handler, body)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("esperava status %d, obteve %d", http.StatusUnprocessableEntity, rec.Code)
		}

		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		if resp["error"] != usecase.ErrInvalidSurname.Error() {
			t.Errorf("esperava erro '%s', obteve '%s'", usecase.ErrInvalidSurname.Error(), resp["error"])
		}
	})

	t.Run("erro 422 ao email ser invalido", func(t *testing.T) {
		mock := newHandlerSuccessMock()
		handler := newTestHandler(mock)

		birthDate := time.Now().AddDate(-20, 0, 0).Format("2006-01-02")
		body := `{"name":"Joao","surname":"Silva","email":"invalido","birthDate":"` + birthDate + `"}`

		rec := executeRegisterRequest(handler, body)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("esperava status %d, obteve %d", http.StatusUnprocessableEntity, rec.Code)
		}

		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		if resp["error"] != usecase.ErrInvalidEmail.Error() {
			t.Errorf("esperava erro '%s', obteve '%s'", usecase.ErrInvalidEmail.Error(), resp["error"])
		}
	})

	t.Run("erro 422 ao usuario ser menor de 18 anos", func(t *testing.T) {
		mock := newHandlerSuccessMock()
		handler := newTestHandler(mock)

		birthDate := time.Now().AddDate(-17, 0, 0).Format("2006-01-02")
		body := `{"name":"Joao","surname":"Silva","email":"joao@email.com","birthDate":"` + birthDate + `"}`

		rec := executeRegisterRequest(handler, body)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("esperava status %d, obteve %d", http.StatusUnprocessableEntity, rec.Code)
		}

		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		if resp["error"] != usecase.ErrUserTooYoung.Error() {
			t.Errorf("esperava erro '%s', obteve '%s'", usecase.ErrUserTooYoung.Error(), resp["error"])
		}
	})

	t.Run("erro 422 ao data de nascimento ser no futuro", func(t *testing.T) {
		mock := newHandlerSuccessMock()
		handler := newTestHandler(mock)

		birthDate := time.Now().AddDate(1, 0, 0).Format("2006-01-02")
		body := `{"name":"Joao","surname":"Silva","email":"joao@email.com","birthDate":"` + birthDate + `"}`

		rec := executeRegisterRequest(handler, body)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("esperava status %d, obteve %d", http.StatusUnprocessableEntity, rec.Code)
		}

		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		if resp["error"] != usecase.ErrFutureBirthDate.Error() {
			t.Errorf("esperava erro '%s', obteve '%s'", usecase.ErrFutureBirthDate.Error(), resp["error"])
		}
	})

	t.Run("erro 400 ao email ja estar em uso", func(t *testing.T) {
		mock := newHandlerSuccessMock()
		mock.getUserByEmailFn = func(email string) (*domain.User, error) {
			return &domain.User{ID: "existing-id", Email: email}, nil
		}
		handler := newTestHandler(mock)

		rec := executeRegisterRequest(handler, validJSON())

		if rec.Code != http.StatusBadRequest {
			t.Errorf("esperava status %d, obteve %d", http.StatusBadRequest, rec.Code)
		}

		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		if resp["error"] != usecase.ErrEmailInUse.Error() {
			t.Errorf("esperava erro '%s', obteve '%s'", usecase.ErrEmailInUse.Error(), resp["error"])
		}
	})

	t.Run("erro 500 ao repositorio falhar", func(t *testing.T) {
		mock := newHandlerSuccessMock()
		mock.saveUserFn = func(user *domain.User) error {
			return errors.New("database connection failed")
		}
		handler := newTestHandler(mock)

		rec := executeRegisterRequest(handler, validJSON())

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("esperava status %d, obteve %d", http.StatusInternalServerError, rec.Code)
		}

		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		if resp["error"] != "Internal server error" {
			t.Errorf("esperava erro 'Internal server error', obteve '%s'", resp["error"])
		}
	})

	t.Run("erro 422 ao data de nascimento ter formato invalido", func(t *testing.T) {
		mock := newHandlerSuccessMock()
		handler := newTestHandler(mock)

		body := `{"name":"Joao","surname":"Silva","email":"joao@email.com","birthDate":"invalido"}`

		rec := executeRegisterRequest(handler, body)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("esperava status %d, obteve %d", http.StatusUnprocessableEntity, rec.Code)
		}

		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		if resp["error"] != usecase.ErrInvalidBirthDate.Error() {
			t.Errorf("esperava erro '%s', obteve '%s'", usecase.ErrInvalidBirthDate.Error(), resp["error"])
		}
	})
}

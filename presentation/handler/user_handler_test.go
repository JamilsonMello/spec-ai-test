package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/example/cadastro-de-usuarios/application/usecase"
	"github.com/example/cadastro-de-usuarios/infrastructure/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestListarUsuarios_Autorizacao(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	listUsersUC := usecase.NewListUsersUseCase(mockRepo)
	userHandler := NewUserHandler(nil, listUsersUC, nil, nil)

	t.Run("FR-1: Autorização via Header - admin role retorna 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/usuarios/listar?page=1&limit=10", nil)
		req.Header.Set("X-User-Role", "admin")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := userHandler.ListUsers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("FR-2: Autorização via Query - admin role retorna 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/usuarios/listar?role=admin", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := userHandler.ListUsers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("FR-3: Acesso Negado - role omitida retorna 403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/usuarios/listar", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := userHandler.ListUsers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, rec.Code)
		assert.Contains(t, rec.Body.String(), "Access denied. Admin role required.")
	})

	t.Run("FR-3: Acesso Negado - role diferente de admin retorna 403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/usuarios/listar", nil)
		req.Header.Set("X-User-Role", "user")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := userHandler.ListUsers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, rec.Code)
		assert.Contains(t, rec.Body.String(), "Access denied. Admin role required.")
	})
}

func TestListarUsuarios_Filtros(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	listUsersUC := usecase.NewListUsersUseCase(mockRepo)
	userHandler := NewUserHandler(nil, listUsersUC, nil, nil)

	t.Run("FR-4: Filtro name ignora espaços e case insensitive", func(t *testing.T) {
		// Build URL with query parameters
		params := url.Values{}
		params.Set("role", "admin")
		params.Set("name", "  JoHn  ")
		req := httptest.NewRequest(http.MethodGet, "/usuarios/listar?"+params.Encode(), nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := userHandler.ListUsers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		// Verify that the filter passed to repository is trimmed and lowercased
		assert.Equal(t, "john", mockRepo.LastListUsersFilter.Name)
	})

	t.Run("Filtro email também é sanitizado", func(t *testing.T) {
		params := url.Values{}
		params.Set("role", "admin")
		params.Set("email", "  EXAMPLE@Domain.COM  ")
		req := httptest.NewRequest(http.MethodGet, "/usuarios/listar?"+params.Encode(), nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := userHandler.ListUsers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "example@domain.com", mockRepo.LastListUsersFilter.Email)
	})
}

func TestListarUsuarios_Paginacao(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	listUsersUC := usecase.NewListUsersUseCase(mockRepo)
	userHandler := NewUserHandler(nil, listUsersUC, nil, nil)

	t.Run("FR-5: Paginação padrão - omissão de page e limit resulta em page=1 e limit=30", func(t *testing.T) {
		params := url.Values{}
		params.Set("role", "admin")
		req := httptest.NewRequest(http.MethodGet, "/usuarios/listar?"+params.Encode(), nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := userHandler.ListUsers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, 1, mockRepo.LastListUsersPage)
		assert.Equal(t, 30, mockRepo.LastListUsersLimit)
	})

	t.Run("FR-6: Limite máximo - limit>30 é automaticamente reduzido para 30", func(t *testing.T) {
		params := url.Values{}
		params.Set("role", "admin")
		params.Set("limit", "100")
		req := httptest.NewRequest(http.MethodGet, "/usuarios/listar?"+params.Encode(), nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := userHandler.ListUsers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, 30, mockRepo.LastListUsersLimit)
	})

	t.Run("Page negativa ou zero é convertida para 1", func(t *testing.T) {
		params := url.Values{}
		params.Set("role", "admin")
		params.Set("page", "0")
		req := httptest.NewRequest(http.MethodGet, "/usuarios/listar?"+params.Encode(), nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := userHandler.ListUsers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, 1, mockRepo.LastListUsersPage)
	})

	t.Run("Limit negativo ou zero é convertido para 30", func(t *testing.T) {
		params := url.Values{}
		params.Set("role", "admin")
		params.Set("limit", "-5")
		req := httptest.NewRequest(http.MethodGet, "/usuarios/listar?"+params.Encode(), nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := userHandler.ListUsers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, 30, mockRepo.LastListUsersLimit)
	})
}

func TestListarUsuarios_ErroRepository(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	mockRepo.ListUsersError = assert.AnError // Simulate repository error
	listUsersUC := usecase.NewListUsersUseCase(mockRepo)
	userHandler := NewUserHandler(nil, listUsersUC, nil, nil)

	t.Run("FR-7: Erro de Banco retorna 500 Internal Server Error", func(t *testing.T) {
		params := url.Values{}
		params.Set("role", "admin")
		req := httptest.NewRequest(http.MethodGet, "/usuarios/listar?"+params.Encode(), nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := userHandler.ListUsers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Internal server error")
	})
}

func TestListarUsuarios_ParametrosInvalidos(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	listUsersUC := usecase.NewListUsersUseCase(mockRepo)
	userHandler := NewUserHandler(nil, listUsersUC, nil, nil)

	t.Run("Parâmetros inválidos - letras em page retorna 400", func(t *testing.T) {
		params := url.Values{}
		params.Set("role", "admin")
		params.Set("page", "abc")
		req := httptest.NewRequest(http.MethodGet, "/usuarios/listar?"+params.Encode(), nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := userHandler.ListUsers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid request")
	})

	t.Run("Parâmetros inválidos - letras em limit retorna 400", func(t *testing.T) {
		params := url.Values{}
		params.Set("role", "admin")
		params.Set("limit", "xyz")
		req := httptest.NewRequest(http.MethodGet, "/usuarios/listar?"+params.Encode(), nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := userHandler.ListUsers(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid request")
	})
}

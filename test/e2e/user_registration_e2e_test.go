package e2e

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/example/cadastro-de-usuarios/application/usecase"
	"github.com/example/cadastro-de-usuarios/domain"
	pkgdb "github.com/example/cadastro-de-usuarios/pkg/db"
	"github.com/example/cadastro-de-usuarios/presentation/handler"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

type RegisterUserRequest struct {
	Nome           string `json:"nome"`
	Sobrenome      string `json:"sobrenome"`
	Email          string `json:"email"`
	DataNascimento string `json:"dataNascimento"`
	Senha          string `json:"senha"`
}

type RegisterUserResponse struct {
	ID             string `json:"id"`
	Nome           string `json:"nome"`
	Sobrenome      string `json:"sobrenome"`
	Email          string `json:"email"`
	DataNascimento string `json:"dataNascimento"`
}

type UserRepositorySpy struct {
	delegate   domain.UserRepository
	forceError error
}

func NewUserRepositorySpy(delegate domain.UserRepository) *UserRepositorySpy {
	return &UserRepositorySpy{delegate: delegate}
}

func (r *UserRepositorySpy) SaveUser(user *domain.User) error {
	if r.forceError != nil {
		return r.forceError
	}
	return r.delegate.SaveUser(user)
}

func (r *UserRepositorySpy) GetUserByEmail(email string) (*domain.User, error) {
	return r.delegate.GetUserByEmail(email)
}

func (r *UserRepositorySpy) GetUserByID(id string) (*domain.User, error) {
	return r.delegate.GetUserByID(id)
}

func (r *UserRepositorySpy) DeleteUser(id string) error {
	return r.delegate.DeleteUser(id)
}

func (r *UserRepositorySpy) UpdateUser(user *domain.User) error {
	return r.delegate.UpdateUser(user)
}

func (r *UserRepositorySpy) ListUsers(filter domain.UserFilter, page int, limit int) ([]*domain.User, int, error) {
	return r.delegate.ListUsers(filter, page, limit)
}

func setupTestServer(db *sql.DB) (*echo.Echo, *UserRepositorySpy) {
	delegateRepo, err := NewPostgreSQLUserRepository(db)
	if err != nil {
		panic(err)
	}
	userRepoSpy := NewUserRepositorySpy(delegateRepo)

	registerUserUC := usecase.NewRegisterUserUseCase(userRepoSpy)
	userHandler := handler.NewUserHandler(registerUserUC, nil, nil, nil)

	e := echo.New()
	e.POST("/api/users/register", userHandler.RegisterUser)

	return e, userRepoSpy
}

func cleanupDatabase(db *sql.DB) {
	db.Exec("DELETE FROM users")
}

func TestUserRegistration_Success(t *testing.T) {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping E2E test")
	}

	db, err := pkgdb.Connect(testDBURL)
	require.NoError(t, err)
	defer db.Close()

	cleanupDatabase(db)
	defer cleanupDatabase(db)

	e, _ := setupTestServer(db)

	reqBody := RegisterUserRequest{
		Nome:           "João",
		Sobrenome:      "Silva",
		Email:          "joao.silva@example.com",
		DataNascimento: "1990-01-01",
		Senha:          "SenhaSegura123",
	}
	reqJSON, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp RegisterUserResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "João", resp.Nome)
	assert.Equal(t, "Silva", resp.Sobrenome)
	assert.Equal(t, "joao.silva@example.com", resp.Email)
	assert.Equal(t, "1990-01-01", resp.DataNascimento)
	assert.NotEmpty(t, resp.ID)

	var persistedName, persistedSurname, persistedEmail, persistedBirthDate, persistedPassword string
	err = db.QueryRow("SELECT name, surname, email, birth_date, password FROM users WHERE email = $1", "joao.silva@example.com").Scan(
		&persistedName, &persistedSurname, &persistedEmail, &persistedBirthDate, &persistedPassword,
	)
	require.NoError(t, err)
	assert.Equal(t, "João", persistedName)
	assert.Equal(t, "Silva", persistedSurname)
	assert.Equal(t, "joao.silva@example.com", persistedEmail)
	assert.Equal(t, "1990-01-01", persistedBirthDate)
	assert.NotEqual(t, "SenhaSegura123", persistedPassword)
}

func TestUserRegistration_DuplicateEmail(t *testing.T) {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping E2E test")
	}

	db, err := pkgdb.Connect(testDBURL)
	require.NoError(t, err)
	defer db.Close()

	cleanupDatabase(db)
	defer cleanupDatabase(db)

	e, _ := setupTestServer(db)

	firstUser := RegisterUserRequest{
		Nome:           "Maria",
		Sobrenome:      "Santos",
		Email:          "maria.santos@example.com",
		DataNascimento: "1985-05-15",
	}
	reqJSON, _ := json.Marshal(firstUser)
	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	duplicateUser := RegisterUserRequest{
		Nome:           "João",
		Sobrenome:      "Silva",
		Email:          "maria.santos@example.com",
		DataNascimento: "1990-01-01",
	}
	reqJSON, _ = json.Marshal(duplicateUser)
	req = httptest.NewRequest(http.MethodPost, "/usuarios", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)

	var errResp ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "Email já em uso", errResp.Error)
	assert.Equal(t, "DUPLICATE_EMAIL", errResp.Code)
}

func TestUserRegistration_InvalidName(t *testing.T) {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping E2E test")
	}

	db, err := pkgdb.Connect(testDBURL)
	require.NoError(t, err)
	defer db.Close()

	cleanupDatabase(db)
	defer cleanupDatabase(db)

	e, _ := setupTestServer(db)

	reqBody := RegisterUserRequest{
		Nome:           "J",
		Sobrenome:      "Silva",
		Email:          "joao.silva@example.com",
		DataNascimento: "1990-01-01",
	}
	reqJSON, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "Nome inválido", errResp.Error)
	assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
}

func TestUserRegistration_InvalidSurname(t *testing.T) {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping E2E test")
	}

	db, err := pkgdb.Connect(testDBURL)
	require.NoError(t, err)
	defer db.Close()

	cleanupDatabase(db)
	defer cleanupDatabase(db)

	e, _ := setupTestServer(db)

	reqBody := RegisterUserRequest{
		Nome:           "João",
		Sobrenome:      "S",
		Email:          "joao.silva@example.com",
		DataNascimento: "1990-01-01",
	}
	reqJSON, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "Sobrenome inválido", errResp.Error)
	assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
}

func TestUserRegistration_InvalidEmailFormat(t *testing.T) {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping E2E test")
	}

	db, err := pkgdb.Connect(testDBURL)
	require.NoError(t, err)
	defer db.Close()

	cleanupDatabase(db)
	defer cleanupDatabase(db)

	e, _ := setupTestServer(db)

	reqBody := RegisterUserRequest{
		Nome:           "João",
		Sobrenome:      "Silva",
		Email:          "invalid-email",
		DataNascimento: "1990-01-01",
	}
	reqJSON, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "Email em formato incorreto", errResp.Error)
	assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
}

func TestUserRegistration_InvalidBirthDateFormat(t *testing.T) {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping E2E test")
	}

	db, err := pkgdb.Connect(testDBURL)
	require.NoError(t, err)
	defer db.Close()

	cleanupDatabase(db)
	defer cleanupDatabase(db)

	e, _ := setupTestServer(db)

	reqBody := RegisterUserRequest{
		Nome:           "João",
		Sobrenome:      "Silva",
		Email:          "joao.silva@example.com",
		DataNascimento: "01-01-1990",
	}
	reqJSON, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "Data de nascimento inválida", errResp.Error)
	assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
}

func TestUserRegistration_FutureBirthDate(t *testing.T) {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping E2E test")
	}

	db, err := pkgdb.Connect(testDBURL)
	require.NoError(t, err)
	defer db.Close()

	cleanupDatabase(db)
	defer cleanupDatabase(db)

	e, _ := setupTestServer(db)

	reqBody := RegisterUserRequest{
		Nome:           "João",
		Sobrenome:      "Silva",
		Email:          "joao.silva@example.com",
		DataNascimento: "2030-01-01",
	}
	reqJSON, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "Data de nascimento inválida", errResp.Error)
	assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
}

func TestUserRegistration_UserUnder18(t *testing.T) {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping E2E test")
	}

	db, err := pkgdb.Connect(testDBURL)
	require.NoError(t, err)
	defer db.Close()

	cleanupDatabase(db)
	defer cleanupDatabase(db)

	e, _ := setupTestServer(db)

	reqBody := RegisterUserRequest{
		Nome:           "João",
		Sobrenome:      "Silva",
		Email:          "joao.silva@example.com",
		DataNascimento: "2015-01-01",
	}
	reqJSON, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "Usuário menor de 18 anos", errResp.Error)
	assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
}

func TestUserRegistration_ShortPassword(t *testing.T) {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping E2E test")
	}

	db, err := pkgdb.Connect(testDBURL)
	require.NoError(t, err)
	defer db.Close()

	cleanupDatabase(db)
	defer cleanupDatabase(db)

	e, _ := setupTestServer(db)

	reqBody := RegisterUserRequest{
		Nome:           "João",
		Sobrenome:      "Silva",
		Email:          "joao.silva@example.com",
		DataNascimento: "1990-01-01",
		Senha:          "abc",
	}
	reqJSON, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "Senha muito curta", errResp.Error)
	assert.Equal(t, "VALIDATION_ERROR", errResp.Code)
}

func TestUserRegistration_PersistenceFailure(t *testing.T) {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping E2E test")
	}

	db, err := pkgdb.Connect(testDBURL)
	require.NoError(t, err)
	defer db.Close()

	cleanupDatabase(db)
	defer cleanupDatabase(db)

	e, userRepoSpy := setupTestServer(db)

	userRepoSpy.forceError = errors.New("database connection failed")

	reqBody := RegisterUserRequest{
		Nome:           "João",
		Sobrenome:      "Silva",
		Email:          "joao.silva@example.com",
		DataNascimento: "1990-01-01",
	}
	reqJSON, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var errResp ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "Erro interno do servidor ao processar o cadastro", errResp.Error)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", errResp.Code)
}

type PostgreSQLUserRepository struct {
	db *sql.DB
}

func NewPostgreSQLUserRepository(db *sql.DB) (*PostgreSQLUserRepository, error) {
	return &PostgreSQLUserRepository{db: db}, nil
}

func (r *PostgreSQLUserRepository) SaveUser(user *domain.User) error {
	query := `INSERT INTO users (id, name, surname, email, birth_date, password, recovery_token, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Surname, user.Email, user.BirthDate, user.Password, user.RecoveryToken, user.Role, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

func (r *PostgreSQLUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	query := `SELECT id, name, surname, email, birth_date, password, recovery_token, role, created_at
		FROM users WHERE email = $1`
	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Surname, &user.Email,
		&user.BirthDate, &user.Password, &user.RecoveryToken, &user.Role, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (r *PostgreSQLUserRepository) GetUserByID(id string) (*domain.User, error) {
	query := `SELECT id, name, surname, email, birth_date, password, recovery_token, role, created_at
		FROM users WHERE id = $1`
	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Surname, &user.Email,
		&user.BirthDate, &user.Password, &user.RecoveryToken, &user.Role, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (r *PostgreSQLUserRepository) DeleteUser(id string) error {
	result, err := r.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *PostgreSQLUserRepository) UpdateUser(user *domain.User) error {
	query := `UPDATE users SET name=$1, surname=$2, email=$3, birth_date=$4, password=$5, recovery_token=$6, role=$7
		WHERE id=$8`
	result, err := r.db.Exec(query, user.Name, user.Surname, user.Email, user.BirthDate, user.Password, user.RecoveryToken, user.Role, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *PostgreSQLUserRepository) ListUsers(filter domain.UserFilter, page int, limit int) ([]*domain.User, int, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	idx := 1

	if filter.Name != "" {
		where += fmt.Sprintf(" AND LOWER(name) LIKE $%d", idx)
		args = append(args, "%"+filter.Name+"%")
		idx++
	}
	if filter.Email != "" {
		where += fmt.Sprintf(" AND LOWER(email) LIKE $%d", idx)
		args = append(args, "%"+filter.Email+"%")
		idx++
	}

	var total int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM users "+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	listArgs := append(args, limit, (page-1)*limit)
	listQuery := fmt.Sprintf(
		`SELECT id, name, surname, email, birth_date, password, recovery_token, role, created_at
		FROM users %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		where, idx, idx+1,
	)

	rows, err := r.db.Query(listQuery, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		user := &domain.User{}
		if err := rows.Scan(
			&user.ID, &user.Name, &user.Surname, &user.Email,
			&user.BirthDate, &user.Password, &user.RecoveryToken, &user.Role, &user.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate users: %w", err)
	}
	return users, total, nil
}

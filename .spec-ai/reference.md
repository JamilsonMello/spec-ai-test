# Reference Implementation — Existing Feature Template

This file contains a COMPLETE existing feature implementation from the codebase.
Use it as your template — match its patterns, naming, structure, error handling, and code style EXACTLY.
Do NOT invent new patterns. Copy the approach below.

## internal/domain/user.go
```go
package domain

import (
	"regexp"
	"time"
)

type User struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Surname       string    `json:"surname"`
	Email         string    `json:"email"`
	BirthDate     time.Time `json:"birthDate"`
	Password      string    `json:"-"`         
	RecoveryToken string    `json:"-"`         
	Role          string    `json:"role"`      
	CreatedAt     time.Time `json:"createdAt"` 
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var nameSurnameRegex = regexp.MustCompile(`^[a-zA-Z\s]+$`)

func (u *User) IsValidName() bool {
	return len(u.Name) >= 2 && len(u.Name) <= 50 && nameSurnameRegex.MatchString(u.Name)
}

func (u *User) IsValidSurname() bool {
	return len(u.Surname) >= 2 && len(u.Surname) <= 50 && nameSurnameRegex.MatchString(u.Surname)
}

func (u *User) IsValidEmailFormat() bool {
	return emailRegex.MatchString(u.Email)
}

func (u *User) IsAdult() bool {
	eighteenYearsAgo := time.Now().AddDate(-18, 0, 0)
	return u.BirthDate.Before(eighteenYearsAgo) || u.BirthDate.Equal(eighteenYearsAgo)
}

func (u *User) IsPastDate() bool {
	return u.BirthDate.Before(time.Now())
}

func (u *User) IsValidPassword(password string) bool {
	return len(password) >= 8
}

```

## internal/application/usecase/update_user_profile.go
```go
package usecase

import (
	"errors"
	"time"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrInvalidNameUpdate      = errors.New("nome deve ter entre 2 e 50 caracteres e conter apenas letras e espaços")
	ErrInvalidBirthDateUpdate = errors.New("data de nascimento inválida")
	ErrFutureBirthDateUpdate  = errors.New("data de nascimento não pode ser no futuro")
	ErrUserNotFoundUpdate     = errors.New("usuário não encontrado")
)

type UpdateUserProfileRequest struct {
	UserID    string `param:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birthDate"`
}

type UpdateUserProfileResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	BirthDate string `json:"birthDate"`
}

type UpdateUserProfileUseCase struct {
	UserRepository domain.UserRepository
}

func NewUpdateUserProfileUseCase(repo domain.UserRepository) *UpdateUserProfileUseCase {
	return &UpdateUserProfileUseCase{
		UserRepository: repo,
	}
}

func (uc *UpdateUserProfileUseCase) Execute(req UpdateUserProfileRequest) (*UpdateUserProfileResponse, error) {
	user, err := uc.getUser(req.UserID)
	if err != nil {
		return nil, err
	}

	birthDate, err := uc.parseBirthDate(req.BirthDate)
	if err != nil {
		return nil, err
	}

	user.Name = req.Name
	user.BirthDate = birthDate

	if err := uc.validateUpdatedUser(user); err != nil {
		return nil, err
	}

	err = uc.UserRepository.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return uc.buildResponse(user), nil
}

func (uc *UpdateUserProfileUseCase) getUser(userID string) (*domain.User, error) {
	if userID == "" {
		return nil, ErrUserNotFoundUpdate
	}

	user, err := uc.UserRepository.FindUserByUuid(userID)
	if err != nil {
		return nil, ErrUserNotFoundUpdate
	}
	return user, nil
}

func (uc *UpdateUserProfileUseCase) parseBirthDate(birthDateStr string) (time.Time, error) {
	birthDate, err := time.Parse("2006-01-02", birthDateStr)
	if err != nil {
		return time.Time{}, ErrInvalidBirthDateUpdate
	}
	return birthDate, nil
}

func (uc *UpdateUserProfileUseCase) validateUpdatedUser(user *domain.User) error {
	if !user.IsValidName() {
		return ErrInvalidNameUpdate
	}

	if !user.IsPastDate() {
		return ErrFutureBirthDateUpdate
	}

	return nil
}

func (uc *UpdateUserProfileUseCase) buildResponse(user *domain.User) *UpdateUserProfileResponse {
	return &UpdateUserProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		BirthDate: user.BirthDate.Format("2006-01-02"),
	}
}

```

## internal/infrastructure/repository/user_repository.go
```go
package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) SaveUser(user *domain.User) error {
	query := `INSERT INTO users (id, name, surname, email, birth_date, password, recovery_token, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Surname, user.Email, user.BirthDate, user.Password, user.RecoveryToken, user.Role, user.CreatedAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return domain.ErrEmailAlreadyExists
		}
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
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

func (r *UserRepository) FindUserByUuid(id string) (*domain.User, error) {
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
		return nil, fmt.Errorf("failed to find user by uuid: %w", err)
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(user *domain.User) error {
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

func (r *UserRepository) DeleteUser(id string) error {
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

func (r *UserRepository) ListUsers(filter domain.UserFilter, page int, limit int) ([]*domain.User, int, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	idx := 1

	if filter.Name != "" {
		where += fmt.Sprintf(" AND LOWER(name) LIKE $%d", idx)
		args = append(args, "%"+strings.ToLower(filter.Name)+"%")
		idx++
	}
	if filter.Email != "" {
		where += fmt.Sprintf(" AND LOWER(email) LIKE $%d", idx)
		args = append(args, "%"+strings.ToLower(filter.Email)+"%")
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

```

## internal/presentation/handler/user_handler.go
```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
)

type UserHandler struct {
	RegisterUserUseCase      *usecase.RegisterUserUseCase
	ListUsersUseCase         *usecase.ListUsersUseCase
	UpdateUserProfileUseCase *usecase.UpdateUserProfileUseCase
	DeleteUserUseCase        *usecase.DeleteUserUseCase
}

func NewUserHandler(registerUC *usecase.RegisterUserUseCase, listUC *usecase.ListUsersUseCase, updateProfileUC *usecase.UpdateUserProfileUseCase, deleteUC *usecase.DeleteUserUseCase) *UserHandler {
	return &UserHandler{
		RegisterUserUseCase:      registerUC,
		ListUsersUseCase:         listUC,
		UpdateUserProfileUseCase: updateProfileUC,
		DeleteUserUseCase:        deleteUC,
	}
}

func (handler *UserHandler) RegisterUser(c echo.Context) error {
	var req usecase.RegisterUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	resp, err := handler.RegisterUserUseCase.Execute(req)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.JSON(http.StatusOK, resp)
}

func (handler *UserHandler) ListUsers(c echo.Context) error {
	userRole := c.Request().Header.Get("X-User-Role")
	if userRole == "" {
		userRole = c.QueryParam("role")
	}

	if userRole != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied. Admin role required."})
	}

	req := usecase.ListUsersRequest{
		Name:  c.QueryParam("name"),
		Email: c.QueryParam("email"),
	}

	page := 1
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	req.Page = page

	limit := 30
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 30 {
				limit = 30
			}
		}
	}
	req.Limit = limit

	resp, err := handler.ListUsersUseCase.Execute(req)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.JSON(http.StatusOK, resp)
}

func (handler *UserHandler) DeleteUser(c echo.Context) error {
	userRole := c.Request().Header.Get("X-User-Role")
	if userRole == "" {
		userRole = c.QueryParam("role")
	}

	userID := c.Param("id")

	err := handler.DeleteUserUseCase.Execute(userID, userRole)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.NoContent(http.StatusNoContent)
}

func (handler *UserHandler) UpdateUserProfile(c echo.Context) error {
	var req usecase.UpdateUserProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	req.UserID = c.Param("id")

	_, err := handler.UpdateUserProfileUseCase.Execute(req)
	if err != nil {
		statusCode, message := MapErrorToHTTP(err)
		return c.JSON(statusCode, map[string]string{"error": message})
	}

	return c.NoContent(http.StatusNoContent)
}

```


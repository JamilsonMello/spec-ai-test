package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/example/cadastro-de-usuarios/infrastructure/repository"
	"github.com/example/cadastro-de-usuarios/application/usecase"
)

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	RegisterUserUseCase *usecase.RegisterUserUseCase
	ListUsersUseCase    *usecase.ListUsersUseCase
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(registerUC *usecase.RegisterUserUseCase, listUC *usecase.ListUsersUseCase) *UserHandler {
	return &UserHandler{
		RegisterUserUseCase: registerUC,
		ListUsersUseCase:    listUC,
	}
}

// RegisterUser handles the POST /usuarios request.
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req usecase.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := h.RegisterUserUseCase.Execute(req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidName) ||
			errors.Is(err, usecase.ErrInvalidSurname) ||
			errors.Is(err, usecase.ErrInvalidEmail) ||
			errors.Is(err, usecase.ErrInvalidBirthDate) ||
			errors.Is(err, usecase.ErrUserTooYoung) ||
			errors.Is(err, usecase.ErrFutureBirthDate) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		} else if errors.Is(err, repository.ErrEmailAlreadyExists) || errors.Is(err, usecase.ErrEmailInUse) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// ListUsers handles the GET /usuarios request.
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Simulate role check from a middleware.
	role := r.Header.Get("X-User-Role")
	if role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	queryParams := r.URL.Query()
	page, _ := strconv.Atoi(queryParams.Get("page"))
	limit, _ := strconv.Atoi(queryParams.Get("limit"))

	params := usecase.ListUsersParams{
		Name:  queryParams.Get("name"),
		Email: queryParams.Get("email"),
		Page:  page,
		Limit: limit,
	}

	users, err := h.ListUsersUseCase.Execute(params)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

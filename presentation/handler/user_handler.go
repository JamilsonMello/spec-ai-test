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

// ListUsers handles the GET /usuarios/listar request.
// This endpoint requires admin role for access.
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Check for admin role (basic auth check via header or query param for demo)
	// In production, this should use proper JWT/session validation
	userRole := r.Header.Get("X-User-Role")
	if userRole == "" {
		userRole = r.URL.Query().Get("role")
	}

	// Only allow admin users to access this endpoint
	if userRole != "admin" {
		http.Error(w, "Access denied. Admin role required.", http.StatusForbidden)
		return
	}

	// Parse query parameters
	query := r.URL.Query()

	// Build request DTO
	req := usecase.ListUsersRequest{
		Name:  query.Get("name"),
		Email: query.Get("email"),
	}

	// Parse page (default to 1)
	page := 1
	if pageStr := query.Get("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	req.Page = page

	// Parse limit (default to 30, max 30)
	limit := 30
	if limitStr := query.Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 30 {
				limit = 30
			}
		}
	}
	req.Limit = limit

	// Execute use case
	resp, err := h.ListUsersUseCase.Execute(req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

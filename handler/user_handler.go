package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/example/cadastro-de-usuarios/infra/repository"
	"github.com/example/cadastro-de-usuarios/usecase"
)

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	RegisterUserUseCase *usecase.RegisterUserUseCase
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(registerUC *usecase.RegisterUserUseCase) *UserHandler {
	return &UserHandler{
		RegisterUserUseCase: registerUC,
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

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/example/cadastro-de-usuarios/application/usecase"
)

// PasswordRecoveryHandler handles HTTP requests related to password recovery.
type PasswordRecoveryHandler struct {
	RequestPasswordRecoveryUseCase *usecase.RequestPasswordRecoveryUseCase
}

// NewPasswordRecoveryHandler creates a new PasswordRecoveryHandler.
func NewPasswordRecoveryHandler(recoveryUC *usecase.RequestPasswordRecoveryUseCase) *PasswordRecoveryHandler {
	return &PasswordRecoveryHandler{
		RequestPasswordRecoveryUseCase: recoveryUC,
	}
}

// RequestPasswordRecovery handles the POST /password-recovery request.
func (h *PasswordRecoveryHandler) RequestPasswordRecovery(w http.ResponseWriter, r *http.Request) {
	var req usecase.RequestPasswordRecoveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := h.RequestPasswordRecoveryUseCase.Execute(req)
	if err != nil {
		if err == usecase.ErrUserNotFound {
			// Return generic message to avoid email enumeration
			http.Error(w, "Se o email existir em nossa base, você receberá instruções de recuperação", http.StatusOK)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestRequestPasswordRecovery_Success(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUser(t, userID, "Recovery", "User", "recovery@test.com")

	body := `{"email":"recovery@test.com"}`
	resp, err := http.Post(testServer.URL+"/password-recovery", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(b))
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["token"] == nil {
		t.Fatal("expected token in response")
	}
	if result["message"] == nil {
		t.Fatal("expected message in response")
	}
}

func TestRequestPasswordRecovery_UserNotFound(t *testing.T) {
	cleanTables(t)

	body := `{"email":"nonexistent@test.com"}`
	resp, err := http.Post(testServer.URL+"/password-recovery", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 (security), got %d", resp.StatusCode)
	}
}

func TestRequestPasswordRecovery_InvalidPayload(t *testing.T) {
	cleanTables(t)

	body := `invalid-json`
	resp, err := http.Post(testServer.URL+"/password-recovery", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestResetPassword_Success(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUser(t, userID, "Reset", "User", "reset@test.com")

	recoverBody := `{"email":"reset@test.com"}`
	recoverResp, err := http.Post(testServer.URL+"/password-recovery", "application/json", strings.NewReader(recoverBody))
	if err != nil {
		t.Fatalf("recovery request failed: %v", err)
	}
	defer recoverResp.Body.Close()

	var recoverResult map[string]interface{}
	json.NewDecoder(recoverResp.Body).Decode(&recoverResult)
	token, ok := recoverResult["token"].(string)
	if !ok || token == "" {
		t.Fatal("expected token from recovery response")
	}

	resetBody := `{"token":"` + token + `","newPassword":"novaSenha123","confirmPassword":"novaSenha123"}`
	resp, err := http.Post(testServer.URL+"/password-recovery/reset", "application/json", strings.NewReader(resetBody))
	if err != nil {
		t.Fatalf("reset request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(b))
	}

	var resetResult map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&resetResult)
	if resetResult["message"] == nil {
		t.Fatal("expected message in response")
	}
}

func TestResetPassword_InvalidToken(t *testing.T) {
	cleanTables(t)

	body := `{"token":"invalid-token-123","newPassword":"novaSenha123","confirmPassword":"novaSenha123"}`
	resp, err := http.Post(testServer.URL+"/password-recovery/reset", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", resp.StatusCode)
	}
}

func TestResetPassword_PasswordMismatch(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUser(t, userID, "Mismatch", "User", "mismatch@test.com")

	recoverBody := `{"email":"mismatch@test.com"}`
	recoverResp, err := http.Post(testServer.URL+"/password-recovery", "application/json", strings.NewReader(recoverBody))
	if err != nil {
		t.Fatalf("recovery request failed: %v", err)
	}
	defer recoverResp.Body.Close()

	var recoverResult map[string]interface{}
	json.NewDecoder(recoverResp.Body).Decode(&recoverResult)
	token := recoverResult["token"].(string)

	resetBody := `{"token":"` + token + `","newPassword":"novaSenha123","confirmPassword":"outraSenha456"}`
	resp, err := http.Post(testServer.URL+"/password-recovery/reset", "application/json", strings.NewReader(resetBody))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", resp.StatusCode)
	}
}

func TestResetPassword_TooShort(t *testing.T) {
	cleanTables(t)

	body := `{"token":"some-token","newPassword":"short","confirmPassword":"short"}`
	resp, err := http.Post(testServer.URL+"/password-recovery/reset", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", resp.StatusCode)
	}
}

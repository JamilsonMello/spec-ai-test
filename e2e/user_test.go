package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestRegisterUser_Success(t *testing.T) {
	cleanTables(t)

	body := `{"name":"Joao","surname":"Silva","email":"joao@test.com","birthDate":"1990-01-15"}`
	resp, err := http.Post(testServer.URL+"/usuarios", "application/json", strings.NewReader(body))
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

	if result["id"] == nil {
		t.Fatal("expected id in response")
	}
	if result["name"] != "Joao" {
		t.Fatalf("expected name Joao, got %v", result["name"])
	}
}

func TestRegisterUser_InvalidEmail(t *testing.T) {
	cleanTables(t)

	body := `{"name":"Joao","surname":"Silva","email":"invalid","birthDate":"1990-01-15"}`
	resp, err := http.Post(testServer.URL+"/usuarios", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", resp.StatusCode)
	}
}

func TestRegisterUser_DuplicateEmail(t *testing.T) {
	cleanTables(t)

	body := `{"name":"Joao","surname":"Silva","email":"dup@test.com","birthDate":"1990-01-15"}`
	resp, err := http.Post(testServer.URL+"/usuarios", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	resp.Body.Close()

	resp, err = http.Post(testServer.URL+"/usuarios", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestRegisterUser_MissingFields(t *testing.T) {
	cleanTables(t)

	body := `{"name":"","surname":"Silva","email":"test@test.com","birthDate":"1990-01-15"}`
	resp, err := http.Post(testServer.URL+"/usuarios", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", resp.StatusCode)
	}
}

func TestRegisterUser_Underage(t *testing.T) {
	cleanTables(t)

	body := `{"name":"Joao","surname":"Silva","email":"young@test.com","birthDate":"2020-01-15"}`
	resp, err := http.Post(testServer.URL+"/usuarios", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", resp.StatusCode)
	}
}

func TestListUsers_Success(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUserWithRole(t, userID, "Admin", "User", "admin@test.com", "admin")

	req, _ := http.NewRequest("GET", testServer.URL+"/usuarios/listar?role=admin", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("X-User-Role", "admin")

	resp, err := http.DefaultClient.Do(req)
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

	if result["users"] == nil {
		t.Fatal("expected users in response")
	}
}

func TestListUsers_Unauthorized(t *testing.T) {
	cleanTables(t)

	req, _ := http.NewRequest("GET", testServer.URL+"/usuarios/listar", nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestListUsers_ForbiddenNonAdmin(t *testing.T) {
	cleanTables(t)

	req, _ := http.NewRequest("GET", testServer.URL+"/usuarios/listar", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("X-User-Role", "user")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}

func TestUpdateUserProfile_Success(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUser(t, userID, "Joao", "Silva", "joao@test.com")

	body := `{"name":"Joao Atualizado","birthDate":"1991-06-20"}`
	req, _ := http.NewRequest("PUT", testServer.URL+"/usuarios/"+userID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 204, got %d: %s", resp.StatusCode, string(b))
	}
}

func TestUpdateUserProfile_NotFound(t *testing.T) {
	cleanTables(t)

	fakeID := uuid.New().String()
	body := `{"name":"Joao Atualizado","birthDate":"1991-06-20"}`
	req, _ := http.NewRequest("PUT", testServer.URL+"/usuarios/"+fakeID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestUpdateUserProfile_InvalidName(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUser(t, userID, "Joao", "Silva", "joao2@test.com")

	body := `{"name":"J","birthDate":"1991-06-20"}`
	req, _ := http.NewRequest("PUT", testServer.URL+"/usuarios/"+userID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", resp.StatusCode)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUser(t, userID, "ToDelete", "User", "delete@test.com")

	req, _ := http.NewRequest("DELETE", testServer.URL+"/usuarios/"+userID, nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("X-User-Role", "admin")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 204, got %d: %s", resp.StatusCode, string(b))
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	cleanTables(t)

	fakeID := uuid.New().String()
	req, _ := http.NewRequest("DELETE", testServer.URL+"/usuarios/"+fakeID, nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("X-User-Role", "admin")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestDeleteUser_ForbiddenNonAdmin(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUser(t, userID, "SomeUser", "Test", "some@test.com")

	req, _ := http.NewRequest("DELETE", testServer.URL+"/usuarios/"+userID, nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("X-User-Role", "user")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}

func TestUploadProfilePicture_Success(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUser(t, userID, "PhotoUser", "Test", "photo@test.com")

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("foto", "profile.jpg")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	part.Write([]byte("fake-jpg-content"))
	partHeader := fmt.Sprintf("image/jpeg")
	_ = partHeader
	writer.Close()

	req, _ := http.NewRequest("POST", testServer.URL+"/usuarios/"+userID+"/foto-perfil", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer valid-token")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(b))
	}
}

func TestUploadProfilePicture_UserNotFound(t *testing.T) {
	cleanTables(t)

	fakeID := uuid.New().String()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("foto", "profile.jpg")
	part.Write([]byte("fake-jpg-content"))
	writer.Close()

	req, _ := http.NewRequest("POST", testServer.URL+"/usuarios/"+fakeID+"/foto-perfil", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer valid-token")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 404, got %d: %s", resp.StatusCode, string(b))
	}
}

func TestUploadProfilePicture_NoAuth(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("foto", "profile.jpg")
	part.Write([]byte("fake-jpg-content"))
	writer.Close()

	req, _ := http.NewRequest("POST", testServer.URL+"/usuarios/"+userID+"/foto-perfil", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

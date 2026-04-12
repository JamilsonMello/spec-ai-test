package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestCreatePost_Success(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUser(t, userID, "PostAuthor", "Test", "author@test.com")
	jwtMock.Subject = userID

	body := `{"content":"Meu primeiro post de teste"}`
	req, _ := http.NewRequest("POST", testServer.URL+"/posts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 201, got %d: %s", resp.StatusCode, string(b))
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["id"] == nil {
		t.Fatal("expected id in response")
	}
	if result["content"] != "Meu primeiro post de teste" {
		t.Fatalf("expected content match, got %v", result["content"])
	}

	jwtMock.Subject = "test-user-id"
}

func TestCreatePost_EmptyContent(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUser(t, userID, "PostAuthor", "Test", "author2@test.com")

	body := `{"content":""}`
	req, _ := http.NewRequest("POST", testServer.URL+"/posts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 422, got %d: %s", resp.StatusCode, string(b))
	}
}

func TestCreatePost_Unauthorized(t *testing.T) {
	cleanTables(t)

	body := `{"content":"Test post"}`
	req, _ := http.NewRequest("POST", testServer.URL+"/posts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestUpdatePost_Success(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUser(t, userID, "PostAuthor", "Test", "author3@test.com")
	jwtMock.Subject = userID

	postID := uuid.New().String()
	_, err := testDB.Exec(
		`INSERT INTO posts (id, content, author_id, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())`,
		postID, "Original content", userID,
	)
	if err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	body := `{"content":"Conteudo atualizado"}`
	req, _ := http.NewRequest("PUT", testServer.URL+"/posts/"+postID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
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

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["content"] != "Conteudo atualizado" {
		t.Fatalf("expected updated content, got %v", result["content"])
	}

	jwtMock.Subject = "test-user-id"
}

func TestUpdatePost_NotFound(t *testing.T) {
	cleanTables(t)

	fakePostID := uuid.New().String()

	body := `{"content":"Updated content"}`
	req, _ := http.NewRequest("PUT", testServer.URL+"/posts/"+fakePostID, strings.NewReader(body))
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

func TestUpdatePost_UnauthorizedNotAuthor(t *testing.T) {
	cleanTables(t)

	authorID := uuid.New().String()
	otherUserID := uuid.New().String()
	createTestUser(t, authorID, "Author", "Test", "author4@test.com")
	createTestUser(t, otherUserID, "Other", "User", "other@test.com")

	jwtMock.Subject = otherUserID

	postID := uuid.New().String()
	_, err := testDB.Exec(
		`INSERT INTO posts (id, content, author_id, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())`,
		postID, "Author content", authorID,
	)
	if err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	body := `{"content":"Trying to update others post"}`
	req, _ := http.NewRequest("PUT", testServer.URL+"/posts/"+postID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}

	jwtMock.Subject = "test-user-id"
}

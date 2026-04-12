package e2e

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func createTestPost(t *testing.T, postID, authorID string) {
	t.Helper()
	_, err := testDB.Exec(
		`INSERT INTO posts (id, content, author_id, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())`,
		postID, "Test post content", authorID,
	)
	if err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}
}

func TestCreateComment_Success(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUser(t, userID, "Commenter", "Test", "commenter@test.com")

	postID := uuid.New().String()
	createTestPost(t, postID, userID)

	jwtMock.Subject = userID

	body := `{"content":"Ótimo post!"}`
	req, _ := http.NewRequest("POST", testServer.URL+"/posts/"+postID+"/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

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
	if result["post_id"] != postID {
		t.Fatalf("expected post_id %s, got %v", postID, result["post_id"])
	}
	if result["user_id"] != userID {
		t.Fatalf("expected user_id %s, got %v", userID, result["user_id"])
	}

	jwtMock.Subject = "test-user-id"
}

func TestCreateComment_PostNotFound(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUser(t, userID, "Commenter", "Test", "commenter2@test.com")
	jwtMock.Subject = userID

	fakePostID := uuid.New().String()

	body := `{"content":"Olá"}`
	req, _ := http.NewRequest("POST", testServer.URL+"/posts/"+fakePostID+"/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
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

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["code"] != "POST_NOT_FOUND" {
		t.Fatalf("expected code POST_NOT_FOUND, got %v", result["code"])
	}

	jwtMock.Subject = "test-user-id"
}

func TestCreateComment_EmptyContent(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUser(t, userID, "Commenter", "Test", "commenter3@test.com")

	postID := uuid.New().String()
	createTestPost(t, postID, userID)

	jwtMock.Subject = userID

	body := `{"content":""}`
	req, _ := http.NewRequest("POST", testServer.URL+"/posts/"+postID+"/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 400, got %d: %s", resp.StatusCode, string(b))
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["code"] != "INVALID_INPUT" {
		t.Fatalf("expected code INVALID_INPUT, got %v", result["code"])
	}

	jwtMock.Subject = "test-user-id"
}

func TestCreateComment_Unauthorized(t *testing.T) {
	cleanTables(t)

	body := `{"content":"Test comment"}`
	req, _ := http.NewRequest("POST", testServer.URL+"/posts/some-id/comments", strings.NewReader(body))
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

func TestCreateComment_HTMLSanitized(t *testing.T) {
	cleanTables(t)

	userID := uuid.New().String()
	createTestUser(t, userID, "Commenter", "Test", "commenter4@test.com")

	postID := uuid.New().String()
	createTestPost(t, postID, userID)

	jwtMock.Subject = userID

	body := `{"content":"<script>alert('xss')</script>"}`
	req, _ := http.NewRequest("POST", testServer.URL+"/posts/"+postID+"/comments", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

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

	content, ok := result["content"].(string)
	if !ok {
		t.Fatal("expected content in response")
	}
	if strings.Contains(content, "<script>") {
		t.Fatalf("expected HTML to be sanitized, got: %s", content)
	}

	jwtMock.Subject = "test-user-id"
}

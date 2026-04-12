package e2e

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (s *E2ESuite) createTestUserAndGetID() string {
	return s.createTestUser("Post Author", "Test", "postauthor@test.com", "1990-01-01")
}

func (s *E2ESuite) createTestPost(userID string) string {
	body := `{"content":"Test post content"}`
	req := s.newRequest(http.MethodPost, "/posts", body)
	token := s.generateToken(userID)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-User-ID", userID)
	rec := s.executeRequest(req)
	s.Require().Equal(http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.Require().NoError(err)

	return resp["id"].(string)
}

func (s *E2ESuite) TestCreatePost_Success() {
	userID := s.createTestUserAndGetID()

	body := `{"content":"Test post content"}`
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/posts", body, userID)
	req.Header.Set("X-User-ID", userID)
	rec := s.executeRequest(req)

	s.Equal(http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.NoError(err)
	s.NotEmpty(resp["id"])
	s.Equal("Test post content", resp["content"])
}

func (s *E2ESuite) TestCreatePost_Unauthorized() {
	body := `{"content":"Test post content"}`
	req := s.newRequest(http.MethodPost, "/posts", body)
	rec := s.executeRequest(req)

	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *E2ESuite) TestCreatePost_MissingUserID() {
	body := `{"content":"Test post content"}`
	req := s.newAuthenticatedRequest(http.MethodPost, "/posts", body)
	rec := s.executeRequest(req)

	s.Equal(http.StatusForbidden, rec.Code)
}

func (s *E2ESuite) TestUpdatePost_Success() {
	userID := s.createTestUserAndGetID()
	postID := s.createTestPost(userID)

	body := `{"content":"Updated content"}`
	req := s.newAuthenticatedRequestWithUserID(http.MethodPut, "/posts/"+postID, body, userID)
	rec := s.executeRequest(req)

	s.Equal(http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal("Updated content", resp["content"])
}

func (s *E2ESuite) TestUpdatePost_Unauthorized() {
	body := `{"content":"Updated content"}`
	req := s.newRequest(http.MethodPut, "/posts/"+uuid.New().String(), body)
	rec := s.executeRequest(req)

	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *E2ESuite) TestUpdatePost_NotFound() {
	body := `{"content":"Updated content"}`
	req := s.newAuthenticatedRequest(http.MethodPut, "/posts/"+uuid.New().String(), body)
	rec := s.executeRequest(req)

	s.Equal(http.StatusNotFound, rec.Code)
}

func (s *E2ESuite) TestUpdatePost_ForbiddenOtherUser() {
	userID := s.createTestUserAndGetID()
	postID := s.createTestPost(userID)

	otherUserID := s.createTestUser("Other User", "Test", "other@test.com", "1990-01-01")
	body := `{"content":"Hijack attempt"}`
	req := s.newAuthenticatedRequestWithUserID(http.MethodPut, "/posts/"+postID, body, otherUserID)
	rec := s.executeRequest(req)

	s.Equal(http.StatusForbidden, rec.Code)
}

package e2e

import (
	"encoding/json"
	"fmt"
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

func (s *E2ESuite) createTestCommunity(userID string, name string) string {
	body := fmt.Sprintf(`{"name":"%s","description":"Test community description"}`, name)
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/comunidades", body, userID)
	rec := s.executeRequest(req)
	s.Require().Equal(http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.Require().NoError(err)

	return resp["id"].(string)
}

func (s *E2ESuite) TestCreatePost_WithCommunity_Success() {
	userID := s.createTestUser("Community Post Author", "Test", "communitypost@test.com", "1990-01-01")
	communityID := s.createTestCommunity(userID, "Post Community")

	body := fmt.Sprintf(`{"content":"Post in community","community_id":"%s"}`, communityID)
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/posts", body, userID)
	req.Header.Set("X-User-ID", userID)
	rec := s.executeRequest(req)

	s.Equal(http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.NoError(err)
	s.NotEmpty(resp["id"])
	s.Equal("Post in community", resp["content"])
	s.Equal(communityID, resp["community_id"])
}

func (s *E2ESuite) TestCreatePost_WithNonExistentCommunity_NotFound() {
	userID := s.createTestUser("NotFound Author", "Test", "notfoundcommunity@test.com", "1990-01-01")
	fakeCommunityID := uuid.New().String()

	body := fmt.Sprintf(`{"content":"Post in fake community","community_id":"%s"}`, fakeCommunityID)
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/posts", body, userID)
	req.Header.Set("X-User-ID", userID)
	rec := s.executeRequest(req)

	s.Equal(http.StatusNotFound, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal("community_not_found", resp["error"])
}

func (s *E2ESuite) TestCreatePost_WithInvalidCommunityUUID_BadRequest() {
	userID := s.createTestUser("Invalid UUID Author", "Test", "invaliduuid@test.com", "1990-01-01")

	body := `{"content":"Post with bad uuid","community_id":"not-a-valid-uuid"}`
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/posts", body, userID)
	req.Header.Set("X-User-ID", userID)
	rec := s.executeRequest(req)

	s.Equal(http.StatusBadRequest, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal("invalid_uuid_format", resp["error"])
}

package e2e

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (s *E2ESuite) createTestComment(userID, postID string) string {
	body := `{"content":"Test comment content"}`
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/posts/"+postID+"/comments", body, userID)
	rec := s.executeRequest(req)
	s.Require().Equal(http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.Require().NoError(err)

	return resp["id"].(string)
}

func (s *E2ESuite) TestCreateComment_Success() {
	userID := s.createTestUserAndGetID()
	postID := s.createTestPost(userID)

	body := `{"content":"This is a comment"}`
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/posts/"+postID+"/comments", body, userID)
	rec := s.executeRequest(req)

	s.Equal(http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.NoError(err)
	s.NotEmpty(resp["id"])
	s.Equal("This is a comment", resp["content"])
}

func (s *E2ESuite) TestCreateComment_Unauthorized() {
	body := `{"content":"This is a comment"}`
	req := s.newRequest(http.MethodPost, "/posts/"+uuid.New().String()+"/comments", body)
	rec := s.executeRequest(req)

	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *E2ESuite) TestCreateComment_InvalidContent() {
	userID := s.createTestUserAndGetID()
	postID := s.createTestPost(userID)

	body := `{"content":""}`
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/posts/"+postID+"/comments", body, userID)
	rec := s.executeRequest(req)

	s.Equal(http.StatusBadRequest, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal("Invalid content", resp["error"])
}

func (s *E2ESuite) TestListComments_Success() {
	userID := s.createTestUserAndGetID()
	postID := s.createTestPost(userID)

	s.createTestComment(userID, postID)
	s.createTestComment(userID, postID)

	req := s.newRequest(http.MethodGet, "/posts/"+postID+"/comments", "")
	rec := s.executeRequest(req)

	s.Equal(http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.NoError(err)
	comments := resp["comments"].([]interface{})
	s.Len(comments, 2)
}

func (s *E2ESuite) TestListComments_EmptyList() {
	userID := s.createTestUserAndGetID()
	postID := s.createTestPost(userID)

	req := s.newRequest(http.MethodGet, "/posts/"+postID+"/comments", "")
	rec := s.executeRequest(req)

	s.Equal(http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.NoError(err)
	comments := resp["comments"].([]interface{})
	s.Len(comments, 0)
}

func (s *E2ESuite) TestToggleReaction_Success() {
	userID := s.createTestUserAndGetID()
	postID := s.createTestPost(userID)
	commentID := s.createTestComment(userID, postID)

	body := `{"type":"LIKE"}`
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/comments/"+commentID+"/reactions", body, userID)
	rec := s.executeRequest(req)

	s.Equal(http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.NoError(err)
	s.NotEmpty(resp["id"])
	s.Equal("LIKE", resp["type"])
}

func (s *E2ESuite) TestToggleReaction_RemoveExisting() {
	userID := s.createTestUserAndGetID()
	postID := s.createTestPost(userID)
	commentID := s.createTestComment(userID, postID)

	body := `{"type":"LIKE"}`
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/comments/"+commentID+"/reactions", body, userID)
	rec := s.executeRequest(req)
	s.Equal(http.StatusCreated, rec.Code)

	req = s.newAuthenticatedRequestWithUserID(http.MethodPost, "/comments/"+commentID+"/reactions", body, userID)
	rec = s.executeRequest(req)
	s.Equal(http.StatusNoContent, rec.Code)
}

func (s *E2ESuite) TestToggleReaction_Unauthorized() {
	body := `{"type":"LIKE"}`
	req := s.newRequest(http.MethodPost, "/comments/"+uuid.New().String()+"/reactions", body)
	rec := s.executeRequest(req)

	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *E2ESuite) TestToggleReaction_CommentNotFound() {
	userID := s.createTestUserAndGetID()

	body := `{"type":"LIKE"}`
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/comments/"+uuid.New().String()+"/reactions", body, userID)
	rec := s.executeRequest(req)

	s.Equal(http.StatusNotFound, rec.Code)
}

func (s *E2ESuite) TestToggleReaction_InvalidType() {
	userID := s.createTestUserAndGetID()
	postID := s.createTestPost(userID)
	commentID := s.createTestComment(userID, postID)

	body := `{"type":"INVALID"}`
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/comments/"+commentID+"/reactions", body, userID)
	rec := s.executeRequest(req)

	s.Equal(http.StatusBadRequest, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal("Invalid reaction type", resp["error"])
}

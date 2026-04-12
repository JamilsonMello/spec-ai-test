package e2e

import (
	"encoding/json"
	"net/http"
)

func (s *E2ESuite) TestCreateCommunity_Success() {
	userID := s.createTestUser("Community Owner", "Test", "communityowner@test.com", "1990-01-01")

	body := `{"name":"Test Community","description":"A test community description"}`
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/comunidades", body, userID)
	rec := s.executeRequest(req)

	s.Equal(http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.NoError(err)
	s.NotEmpty(resp["id"])
	s.Equal("Test Community", resp["name"])
	s.Equal(userID, resp["owner_id"])
}

func (s *E2ESuite) TestCreateCommunity_DuplicateName() {
	userID := s.createTestUser("Community Owner", "Test", "communityowner2@test.com", "1990-01-01")

	body := `{"name":"Duplicate Community","description":"First community"}`
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/comunidades", body, userID)
	rec := s.executeRequest(req)
	s.Equal(http.StatusCreated, rec.Code)

	req2 := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/comunidades", body, userID)
	rec2 := s.executeRequest(req2)
	s.Equal(http.StatusConflict, rec2.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec2.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal("community already exists", resp["error"])
}

func (s *E2ESuite) TestCreateCommunity_Unauthorized() {
	body := `{"name":"Test Community","description":"A test community description"}`
	req := s.newRequest(http.MethodPost, "/comunidades", body)
	rec := s.executeRequest(req)

	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *E2ESuite) TestCreateCommunity_InvalidName() {
	userID := s.createTestUser("Community Owner", "Test", "communityowner3@test.com", "1990-01-01")

	body := `{"name":"","description":"A valid description"}`
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/comunidades", body, userID)
	rec := s.executeRequest(req)

	s.Equal(http.StatusBadRequest, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal("INVALID_INPUT", resp["code"])
}

func (s *E2ESuite) TestCreateCommunity_InvalidDescription() {
	userID := s.createTestUser("Community Owner", "Test", "communityowner4@test.com", "1990-01-01")

	body := `{"name":"Valid Name","description":""}`
	req := s.newAuthenticatedRequestWithUserID(http.MethodPost, "/comunidades", body, userID)
	rec := s.executeRequest(req)

	s.Equal(http.StatusBadRequest, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal("INVALID_INPUT", resp["code"])
}

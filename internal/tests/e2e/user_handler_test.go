package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (s *E2ESuite) createTestUser(name, surname, email, birthDate string) string {
	body := fmt.Sprintf(`{"name":"%s","surname":"%s","email":"%s","birthDate":"%s"}`, name, surname, email, birthDate)
	req := s.newRequest(http.MethodPost, "/usuarios", body)
	rec := s.executeRequest(req)
	s.Require().Equal(http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.Require().NoError(err)

	return resp["id"].(string)
}

func (s *E2ESuite) TestRegisterUser_Success() {
	body := `{"name":"Test User","surname":"Silva","email":"test@test.com","birthDate":"2000-01-01"}`
	req := s.newRequest(http.MethodPost, "/usuarios", body)
	rec := s.executeRequest(req)

	s.Equal(http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	s.NoError(err)
	s.NotEmpty(resp["id"])
	s.Equal("Test User", resp["name"])
	s.Equal("test@test.com", resp["email"])
}

func (s *E2ESuite) TestRegisterUser_InvalidPayload() {
	body := `{"name":"","surname":"","email":"invalid","birthDate":"2000-01-01"}`
	req := s.newRequest(http.MethodPost, "/usuarios", body)
	rec := s.executeRequest(req)

	s.Equal(http.StatusUnprocessableEntity, rec.Code)
}

func (s *E2ESuite) TestRegisterUser_DuplicateEmail() {
	body := `{"name":"Test User","surname":"Silva","email":"dup@test.com","birthDate":"2000-01-01"}`
	req := s.newRequest(http.MethodPost, "/usuarios", body)
	rec := s.executeRequest(req)
	s.Equal(http.StatusOK, rec.Code)

	req = s.newRequest(http.MethodPost, "/usuarios", body)
	rec = s.executeRequest(req)
	s.Equal(http.StatusBadRequest, rec.Code)
}

func (s *E2ESuite) TestListUsers_Success() {
	s.createTestUser("Admin User", "Admin", "admin@test.com", "1990-01-01")

	req := s.newAuthenticatedRequest(http.MethodGet, "/usuarios/listar", "")
	req.Header.Set("X-User-Role", "admin")
	rec := s.executeRequest(req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *E2ESuite) TestListUsers_Unauthorized() {
	req := s.newRequest(http.MethodGet, "/usuarios/listar", "")
	rec := s.executeRequest(req)

	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *E2ESuite) TestListUsers_Forbidden() {
	req := s.newAuthenticatedRequest(http.MethodGet, "/usuarios/listar", "")
	rec := s.executeRequest(req)

	s.Equal(http.StatusForbidden, rec.Code)
}

func (s *E2ESuite) TestDeleteUser_Success() {
	userID := s.createTestUser("Delete User", "Test", "delete@test.com", "1990-01-01")

	req := s.newAuthenticatedRequest(http.MethodDelete, "/usuarios/"+userID, "")
	req.Header.Set("X-User-Role", "admin")
	rec := s.executeRequest(req)

	s.Equal(http.StatusNoContent, rec.Code)
}

func (s *E2ESuite) TestDeleteUser_Unauthorized() {
	req := s.newRequest(http.MethodDelete, "/usuarios/"+uuid.New().String(), "")
	rec := s.executeRequest(req)

	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *E2ESuite) TestDeleteUser_Forbidden() {
	userID := s.createTestUser("Forbidden User", "Test", "forbidden@test.com", "1990-01-01")

	req := s.newAuthenticatedRequest(http.MethodDelete, "/usuarios/"+userID, "")
	rec := s.executeRequest(req)

	s.Equal(http.StatusForbidden, rec.Code)
}

func (s *E2ESuite) TestUpdateUser_Success() {
	userID := s.createTestUser("Old Name", "Surname", "update@test.com", "1990-01-01")

	body := `{"name":"New Name","birthDate":"1990-01-01"}`
	req := s.newAuthenticatedRequest(http.MethodPut, "/usuarios/"+userID, body)
	rec := s.executeRequest(req)

	s.Equal(http.StatusNoContent, rec.Code)
}

func (s *E2ESuite) TestUpdateUser_Unauthorized() {
	body := `{"name":"New Name","birthDate":"1990-01-01"}`
	req := s.newRequest(http.MethodPut, "/usuarios/"+uuid.New().String(), body)
	rec := s.executeRequest(req)

	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *E2ESuite) TestUpdateUser_NotFound() {
	body := `{"name":"New Name","birthDate":"1990-01-01"}`
	req := s.newAuthenticatedRequest(http.MethodPut, "/usuarios/"+uuid.New().String(), body)
	rec := s.executeRequest(req)

	s.Equal(http.StatusNotFound, rec.Code)
}

func (s *E2ESuite) TestUploadProfilePicture_Unauthorized() {
	req := s.newRequest(http.MethodPost, "/usuarios/"+uuid.New().String()+"/foto-perfil", "")
	rec := s.executeRequest(req)

	s.Equal(http.StatusUnauthorized, rec.Code)
}

func (s *E2ESuite) TestPasswordRecovery_Success() {
	s.createTestUser("Recovery User", "Test", "recovery@test.com", "1990-01-01")

	body := `{"email":"recovery@test.com"}`
	req := s.newRequest(http.MethodPost, "/password-recovery", body)
	rec := s.executeRequest(req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *E2ESuite) TestPasswordRecovery_NonExistentEmail() {
	body := `{"email":"nonexistent@test.com"}`
	req := s.newRequest(http.MethodPost, "/password-recovery", body)
	rec := s.executeRequest(req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *E2ESuite) TestResetPassword_InvalidToken() {
	body := `{"token":"invalid-token","newPassword":"newPassword123","confirmPassword":"newPassword123"}`
	req := s.newRequest(http.MethodPost, "/password-recovery/reset", body)
	rec := s.executeRequest(req)

	s.True(rec.Code == http.StatusNotFound || rec.Code == http.StatusUnprocessableEntity)
}

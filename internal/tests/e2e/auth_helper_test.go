package e2e

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/google/uuid"
)

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type jwtClaims struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
}

func (s *E2ESuite) generateToken(userID string) string {
	header := jwtHeader{Alg: "HS256", Typ: "JWT"}
	headerBytes, _ := json.Marshal(header)
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)

	claims := jwtClaims{
		Sub: userID,
		Exp: time.Now().Add(1 * time.Hour).Unix(),
	}
	claimsBytes, _ := json.Marshal(claims)
	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsBytes)

	message := headerEncoded + "." + claimsEncoded
	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write([]byte(message))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return message + "." + signature
}

func (s *E2ESuite) generateExpiredToken(userID string) string {
	header := jwtHeader{Alg: "HS256", Typ: "JWT"}
	headerBytes, _ := json.Marshal(header)
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)

	claims := jwtClaims{
		Sub: userID,
		Exp: time.Now().Add(-1 * time.Hour).Unix(),
	}
	claimsBytes, _ := json.Marshal(claims)
	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsBytes)

	message := headerEncoded + "." + claimsEncoded
	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write([]byte(message))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return message + "." + signature
}

func (s *E2ESuite) newRequest(method, path string, body string) *http.Request {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func (s *E2ESuite) newAuthenticatedRequest(method, path string, body string) *http.Request {
	req := s.newRequest(method, path, body)
	userID := uuid.New().String()
	token := s.generateToken(userID)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func (s *E2ESuite) newAuthenticatedRequestWithUserID(method, path string, body string, userID string) *http.Request {
	req := s.newRequest(method, path, body)
	token := s.generateToken(userID)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func (s *E2ESuite) executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	return rec
}

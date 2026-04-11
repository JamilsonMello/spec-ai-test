package e2e

import (
	"database/sql"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/example/cadastro-de-usuarios/internal/infrastructure"
	"github.com/example/cadastro-de-usuarios/internal/mocks"
	"github.com/example/cadastro-de-usuarios/internal/server"
)

var (
	testDB     *sql.DB
	testServer *httptest.Server
	emailMock  *mocks.EmailSenderMock
	jwtMock    *mocks.JWTValidatorMock
	uploadPath string
)

func TestMain(m *testing.M) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		os.Exit(0)
	}

	db, err := infrastructure.Connect(dsn)
	if err != nil {
		panic(err)
	}
	testDB = db

	emailMock = mocks.NewEmailSenderMock()
	jwtMock = mocks.NewJWTValidatorMock("test-user-id")

	tmpDir, err := os.MkdirTemp("", "e2e-uploads-*")
	if err != nil {
		panic(err)
	}
	uploadPath = tmpDir

	deps := server.Dependencies{
		DB:           testDB,
		EmailSender:  emailMock,
		JWTValidator: jwtMock,
		UploadPath:   uploadPath,
	}

	e := server.SetupServer(deps)
	testServer = httptest.NewServer(e)

	code := m.Run()

	testServer.Close()
	testDB.Close()
	os.RemoveAll(uploadPath)
	os.Exit(code)
}

func cleanTables(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec("DELETE FROM password_recoveries")
	if err != nil {
		t.Fatalf("failed to clean password_recoveries: %v", err)
	}
	_, err = testDB.Exec("DELETE FROM posts")
	if err != nil {
		t.Fatalf("failed to clean posts: %v", err)
	}
	_, err = testDB.Exec("DELETE FROM users")
	if err != nil {
		t.Fatalf("failed to clean users: %v", err)
	}
}

func createTestUser(t *testing.T, id, name, surname, email string) {
	t.Helper()
	_, err := testDB.Exec(
		`INSERT INTO users (id, name, surname, email, birth_date, password, recovery_token, role, created_at)
		VALUES ($1, $2, $3, $4, '1990-01-01', 'hashedpassword', '', 'user', NOW())`,
		id, name, surname, email,
	)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
}

func createTestUserWithRole(t *testing.T, id, name, surname, email, role string) {
	t.Helper()
	_, err := testDB.Exec(
		`INSERT INTO users (id, name, surname, email, birth_date, password, recovery_token, role, created_at)
		VALUES ($1, $2, $3, $4, '1990-01-01', 'hashedpassword', '', $5, NOW())`,
		id, name, surname, email, role,
	)
	if err != nil {
		t.Fatalf("failed to create test user with role: %v", err)
	}
}

package e2e

import (
	"database/sql"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
	"github.com/example/cadastro-de-usuarios/internal/infrastructure"
	"github.com/example/cadastro-de-usuarios/internal/infrastructure/repository"
	"github.com/example/cadastro-de-usuarios/internal/infrastructure/service"
	"github.com/example/cadastro-de-usuarios/internal/infrastructure/storage"
	"github.com/example/cadastro-de-usuarios/internal/presentation/handler"
	"github.com/example/cadastro-de-usuarios/internal/presentation/middleware"
)

type E2ESuite struct {
	suite.Suite
	db     *sql.DB
	echo   *echo.Echo
	secret string
}

func (s *E2ESuite) SetupSuite() {
	dsn := os.Getenv("DATABASE_URL")
	s.Require().NotEmpty(dsn, "DATABASE_URL must be set")

	s.secret = os.Getenv("JWT_SECRET")
	if s.secret == "" {
		s.secret = "test-secret"
		os.Setenv("JWT_SECRET", s.secret)
	}

	db, err := infrastructure.Connect(dsn)
	s.Require().NoError(err)
	s.db = db

	s.setupServer()
}

func (s *E2ESuite) setupServer() {
	userRepo := repository.NewUserRepository(s.db)
	passwordRecoveryRepo := repository.NewPasswordRecoveryRepository(s.db)
	postRepo := repository.NewPostRepository(s.db)
	commentRepo := repository.NewCommentRepository(s.db)
	reactionRepo := repository.NewReactionRepository(s.db)

	emailSender := service.NewEmailSender()
	jwtValidatorService := service.NewJWTValidatorService()

	registerUserUC := usecase.NewRegisterUserUseCase(userRepo)
	listUsersUC := usecase.NewListUsersUseCase(userRepo)
	deleteUserUC := usecase.NewDeleteUserUseCase(userRepo)
	updateUserProfileUC := usecase.NewUpdateUserProfileUseCase(userRepo)
	localStorage := storage.NewLocalStorage("./uploads")
	uploadProfilePictureUC := usecase.NewUploadProfilePictureUseCase(userRepo, localStorage)

	requestPasswordRecoveryUC := usecase.NewRequestPasswordRecoveryUseCase(userRepo, passwordRecoveryRepo, emailSender)
	resetPasswordUC := usecase.NewResetPasswordUseCase(userRepo, passwordRecoveryRepo)
	createPostUC := usecase.NewCreatePostUseCase(postRepo)
	updatePostUC := usecase.NewUpdatePostUseCase(postRepo)
	createCommentUC := usecase.NewCreateCommentUseCase(commentRepo)
	listCommentsUC := usecase.NewListCommentsUseCase(commentRepo)
	toggleReactionUC := usecase.NewToggleCommentReactionUseCase(reactionRepo, commentRepo)
	validateTokenUC := usecase.NewValidateTokenUseCase(jwtValidatorService)

	userHandler := handler.NewUserHandler(registerUserUC, listUsersUC, updateUserProfileUC, deleteUserUC, uploadProfilePictureUC)
	passwordRecoveryHandler := handler.NewPasswordRecoveryHandler(requestPasswordRecoveryUC, resetPasswordUC)
	postHandler := handler.NewPostHandler(createPostUC, updatePostUC)
	commentHandler := handler.NewCommentHandler(createCommentUC, listCommentsUC)
	reactionHandler := handler.NewReactionHandler(toggleReactionUC)

	e := echo.New()

	e.Static("/uploads", "./uploads")

	e.POST("/usuarios", userHandler.RegisterUser)
	e.POST("/password-recovery", passwordRecoveryHandler.RequestPasswordRecovery)
	e.POST("/password-recovery/reset", passwordRecoveryHandler.ResetPassword)

	protected := e.Group("")
	protected.Use(middleware.AuthMiddleware(validateTokenUC))

	protected.GET("/usuarios/listar", userHandler.ListUsers)
	protected.DELETE("/usuarios/:id", userHandler.DeleteUser)
	protected.PUT("/usuarios/:id", userHandler.UpdateUserProfile)
	protected.POST("/usuarios/:id/foto-perfil", userHandler.UploadProfilePicture)
	protected.POST("/posts", postHandler.CreatePost)
	protected.PUT("/posts/:id", postHandler.UpdatePost)
	protected.POST("/posts/:id/comments", commentHandler.CreateComment)
	protected.POST("/comments/:id/reactions", reactionHandler.ToggleReaction)

	e.GET("/posts/:id/comments", commentHandler.ListComments)

	s.echo = e
}

func (s *E2ESuite) TearDownTest() {
	tables := []string{"comment_reactions", "comments", "posts", "password_recoveries", "users"}
	for _, table := range tables {
		_, err := s.db.Exec("TRUNCATE TABLE " + table + " CASCADE")
		s.Require().NoError(err)
	}
}

func (s *E2ESuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
}

func TestE2ESuite(t *testing.T) {
	suite.Run(t, new(E2ESuite))
}

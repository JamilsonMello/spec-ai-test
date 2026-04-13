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
	if dsn == "" {
		s.T().Skip("DATABASE_URL must be set")
	}

	s.secret = os.Getenv("JWT_SECRET")
	if s.secret == "" {
		s.secret = "test-secret"
		os.Setenv("JWT_SECRET", s.secret)
	}

	db, err := infrastructure.Connect(dsn)
	s.Require().NoError(err)

	err = infrastructure.RunMigrations(db)
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
	communityRepo := repository.NewCommunityRepository(s.db)
	productRepo := repository.NewProductRepository(s.db)

	emailSender := service.NewEmailSender()
	bcryptHasher := service.NewBcryptHasher()
	txManager := infrastructure.NewSQLTransactionManager(s.db)
	jwtValidatorService := service.NewJWTValidatorService()

	registerUserUC := usecase.NewRegisterUserUseCase(userRepo)
	listUsersUC := usecase.NewListUsersUseCase(userRepo)
	deleteUserUC := usecase.NewDeleteUserUseCase(userRepo)
	updateUserProfileUC := usecase.NewUpdateUserProfileUseCase(userRepo)
	localStorage := storage.NewLocalStorage("./uploads")
	uploadProfilePictureUC := usecase.NewUploadProfilePictureUseCase(userRepo, localStorage)

	requestPasswordRecoveryUC := usecase.NewRequestPasswordRecoveryUseCase(userRepo, passwordRecoveryRepo, emailSender)
	resetPasswordUC := usecase.NewResetPasswordUseCase(userRepo, userRepo, passwordRecoveryRepo, passwordRecoveryRepo, bcryptHasher, txManager)
	createPostUC := usecase.NewCreatePostUseCase(postRepo, communityRepo)
	updatePostUC := usecase.NewUpdatePostUseCase(postRepo)
	createCommentUC := usecase.NewCreateCommentUseCase(commentRepo)
	listCommentsUC := usecase.NewListCommentsUseCase(commentRepo)
	toggleReactionUC := usecase.NewToggleCommentReactionUseCase(reactionRepo, commentRepo)
	createCommunityUC := usecase.NewCreateCommunityUseCase(communityRepo)
	createProductUC := usecase.NewCreateProductUseCase(productRepo, communityRepo)
	updateProductUC := usecase.NewUpdateProductUseCase(productRepo)
	deleteProductUC := usecase.NewDeleteProductUseCase(productRepo)
	validateTokenUC := usecase.NewValidateTokenUseCase(jwtValidatorService)

	registerUserHandler := handler.NewRegisterUserHandler(registerUserUC)
	listUsersHandler := handler.NewListUsersHandler(listUsersUC)
	deleteUserHandler := handler.NewDeleteUserHandler(deleteUserUC)
	updateUserProfileHandler := handler.NewUpdateUserProfileHandler(updateUserProfileUC)
	uploadProfilePictureHandler := handler.NewUploadProfilePictureHandler(uploadProfilePictureUC)
	passwordRecoveryHandler := handler.NewPasswordRecoveryHandler(requestPasswordRecoveryUC, resetPasswordUC)
	postHandler := handler.NewPostHandler(createPostUC, updatePostUC)
	commentHandler := handler.NewCommentHandler(createCommentUC, listCommentsUC)
	reactionHandler := handler.NewReactionHandler(toggleReactionUC)
	communityHandler := handler.NewCommunityHandler(createCommunityUC)
	productHandler := handler.NewProductHandler(createProductUC, updateProductUC, deleteProductUC)

	e := echo.New()

	e.Static("/uploads", "./uploads")

	e.POST("/users", registerUserHandler.RegisterUser)
	e.POST("/password-recovery", passwordRecoveryHandler.RequestPasswordRecovery)
	e.POST("/password-recovery/reset", passwordRecoveryHandler.ResetPassword)

	protected := e.Group("")
	protected.Use(middleware.AuthMiddleware(validateTokenUC))

	protected.GET("/users", listUsersHandler.ListUsers)
	protected.DELETE("/users/:id", deleteUserHandler.DeleteUser)
	protected.PUT("/users/:id", updateUserProfileHandler.UpdateUserProfile)
	protected.POST("/users/:id/picture", uploadProfilePictureHandler.UploadProfilePicture)
	protected.POST("/posts", postHandler.CreatePost)
	protected.PUT("/posts/:id", postHandler.UpdatePost)
	protected.POST("/posts/:id/comments", commentHandler.CreateComment)
	protected.POST("/comments/:id/reactions", reactionHandler.ToggleReaction)
	protected.POST("/comunidades", communityHandler.CreateCommunity)
	protected.POST("/products", productHandler.CreateProduct)
	protected.PUT("/products/:id", productHandler.UpdateProduct)
	protected.DELETE("/products/:id", productHandler.DeleteProduct)

	e.GET("/posts/:id/comments", commentHandler.ListComments)

	s.echo = e
}

func (s *E2ESuite) TearDownTest() {
	tables := []string{"products", "comment_reactions", "comments", "posts", "communities", "password_recoveries", "users"}
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

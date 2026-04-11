package server

import (
	"database/sql"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
	"github.com/example/cadastro-de-usuarios/internal/domain"
	"github.com/example/cadastro-de-usuarios/internal/infrastructure/repository"
	"github.com/example/cadastro-de-usuarios/internal/presentation/handler"
	"github.com/example/cadastro-de-usuarios/internal/presentation/middleware"
)

type Dependencies struct {
	DB           *sql.DB
	EmailSender  domain.EmailSender
	JWTValidator usecase.TokenValidatorService
	UploadPath   string
}

func SetupServer(deps Dependencies) *echo.Echo {
	userRepo := repository.NewUserRepository(deps.DB)
	passwordRecoveryRepo := repository.NewPasswordRecoveryRepository(deps.DB)
	postRepo := repository.NewPostRepository(deps.DB)

	registerUserUC := usecase.NewRegisterUserUseCase(userRepo)
	listUsersUC := usecase.NewListUsersUseCase(userRepo)
	deleteUserUC := usecase.NewDeleteUserUseCase(userRepo)
	updateUserProfileUC := usecase.NewUpdateUserProfileUseCase(userRepo)
	uploadProfilePictureUC := usecase.NewUploadProfilePictureUseCase(userRepo, deps.UploadPath)

	requestPasswordRecoveryUC := usecase.NewRequestPasswordRecoveryUseCase(userRepo, passwordRecoveryRepo, deps.EmailSender)
	resetPasswordUC := usecase.NewResetPasswordUseCase(userRepo, passwordRecoveryRepo)
	createPostUC := usecase.NewCreatePostUseCase(postRepo)
	updatePostUC := usecase.NewUpdatePostUseCase(postRepo)
	validateTokenUC := usecase.NewValidateTokenUseCase(deps.JWTValidator)

	userHandler := handler.NewUserHandler(registerUserUC, listUsersUC, updateUserProfileUC, deleteUserUC, uploadProfilePictureUC)
	passwordRecoveryHandler := handler.NewPasswordRecoveryHandler(requestPasswordRecoveryUC, resetPasswordUC)
	postHandler := handler.NewPostHandler(createPostUC, updatePostUC)

	e := echo.New()

	e.Static("/uploads", deps.UploadPath)

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

	return e
}

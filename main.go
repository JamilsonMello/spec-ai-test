package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/example/cadastro-de-usuarios/presentation/handler"
	"github.com/example/cadastro-de-usuarios/presentation/middleware"
	"github.com/example/cadastro-de-usuarios/infrastructure/repository"
	"github.com/example/cadastro-de-usuarios/infrastructure/service"
	"github.com/example/cadastro-de-usuarios/application/usecase"
)

func main() {
	// Initialize repositories (Infrastructure layer)
	userRepo := repository.NewInMemoryUserRepository()
	passwordRecoveryRepo := repository.NewInMemoryPasswordRecoveryRepository()
	postRepo := repository.NewInMemoryPostRepository()

	// Initialize services (Infrastructure layer)
	emailService := service.NewEmailService()
	jwtValidatorService := service.NewJWTValidatorService()

	// Initialize use cases (Application layer)
	registerUserUC := usecase.NewRegisterUserUseCase(userRepo)
	listUsersUC := usecase.NewListUsersUseCase(userRepo)
	deleteUserUC := usecase.NewDeleteUserUseCase(userRepo)
	updateUserProfileUC := usecase.NewUpdateUserProfileUseCase(userRepo)

	requestPasswordRecoveryUC := usecase.NewRequestPasswordRecoveryUseCase(userRepo, passwordRecoveryRepo, emailService)
	resetPasswordUC := usecase.NewResetPasswordUseCase(userRepo, passwordRecoveryRepo)
	createPostUC := usecase.NewCreatePostUseCase(postRepo)
	validateTokenUC := usecase.NewValidateTokenUseCase(jwtValidatorService)

	userHandler := handler.NewUserHandler(registerUserUC, listUsersUC, updateUserProfileUC, deleteUserUC)
	passwordRecoveryHandler := handler.NewPasswordRecoveryHandler(requestPasswordRecoveryUC, resetPasswordUC)
	postHandler := handler.NewPostHandler(createPostUC)

	// Set up Echo router
	e := echo.New()

	// Register public routes
	e.POST("/usuarios", userHandler.RegisterUser)
	e.POST("/password-recovery", passwordRecoveryHandler.RequestPasswordRecovery)
	e.POST("/password-recovery/reset", passwordRecoveryHandler.ResetPassword)

	// Register protected routes
	protected := e.Group("")
	protected.Use(middleware.AuthMiddleware(validateTokenUC))

	protected.GET("/usuarios/listar", userHandler.ListUsers)
	protected.DELETE("/usuarios/:id", userHandler.DeleteUser)
	protected.PUT("/usuarios/:id", userHandler.UpdateUserProfile)
	protected.POST("/posts", postHandler.CreatePost)

	// Start the HTTP server
	port := ":8080"
	log.Printf("Server listening on port %s\n", port)
	log.Fatal(e.Start(port))
}

package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/example/cadastro-de-usuarios/presentation/handler"
	"github.com/example/cadastro-de-usuarios/infrastructure/repository"
	"github.com/example/cadastro-de-usuarios/infrastructure/service"
	"github.com/example/cadastro-de-usuarios/application/usecase"
)

func main() {
	// Initialize repositories (Infrastructure layer)
	userRepo := repository.NewInMemoryUserRepository()
	passwordRecoveryRepo := repository.NewInMemoryPasswordRecoveryRepository()

	// Initialize services (Infrastructure layer)
	emailService := service.NewEmailService()

	// Initialize use cases (Application layer)
	registerUserUC := usecase.NewRegisterUserUseCase(userRepo)
	listUsersUC := usecase.NewListUsersUseCase(userRepo)
	requestPasswordRecoveryUC := usecase.NewRequestPasswordRecoveryUseCase(userRepo, passwordRecoveryRepo, emailService)
	resetPasswordUC := usecase.NewResetPasswordUseCase(userRepo, passwordRecoveryRepo)

	// Initialize handlers (Presentation layer)
	userHandler := handler.NewUserHandler(registerUserUC, listUsersUC)
	passwordRecoveryHandler := handler.NewPasswordRecoveryHandler(requestPasswordRecoveryUC, resetPasswordUC)

	// Set up Echo router
	e := echo.New()

	// Register routes
	e.POST("/usuarios", userHandler.RegisterUser)
	e.GET("/usuarios/listar", userHandler.ListUsers)
	e.POST("/password-recovery", passwordRecoveryHandler.RequestPasswordRecovery)
	e.POST("/password-recovery/reset", passwordRecoveryHandler.ResetPassword)

	// Start the HTTP server
	port := ":8080"
	log.Printf("Server listening on port %s\n", port)
	log.Fatal(e.Start(port))
}

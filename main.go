package main

import (
	"log"
	"net/http"

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

	// Initialize handlers (Presentation layer)
	userHandler := handler.NewUserHandler(registerUserUC, listUsersUC)
	passwordRecoveryHandler := handler.NewPasswordRecoveryHandler(requestPasswordRecoveryUC)

	// Set up routes
	http.HandleFunc("/usuarios", userHandler.RegisterUser)
	http.HandleFunc("/usuarios/listar", userHandler.ListUsers)
	http.HandleFunc("/password-recovery", passwordRecoveryHandler.RequestPasswordRecovery)

	// Start the HTTP server
	port := ":8080"
	log.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

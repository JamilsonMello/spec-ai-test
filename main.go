package main

import (
	"log"
	"net/http"

	"github.com/example/cadastro-de-usuarios/handler"
	"github.com/example/cadastro-de-usuarios/infra/repository"
	"github.com/example/cadastro-de-usuarios/usecase"
)

func main() {
	// Initialize dependencies
	userRepo := repository.NewInMemoryUserRepository()
	registerUserUC := usecase.NewRegisterUserUseCase(userRepo)
	listUsersUC := usecase.NewListUsersUseCase(userRepo)
	userHandler := handler.NewUserHandler(registerUserUC, listUsersUC)

	// Set up routes
	http.HandleFunc("/usuarios", userHandler.RegisterUser)
	http.HandleFunc("/usuarios/listar", userHandler.ListUsers)

	// Start the HTTP server
	port := ":8080"
	log.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

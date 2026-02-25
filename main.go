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
	userHandler := handler.NewUserHandler(registerUserUC)

	// Set up routes
	http.HandleFunc("/usuarios", userHandler.RegisterUser)

	// Start the HTTP server
	port := ":8080"
	log.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/example/cadastro-de-usuarios/internal/infrastructure"
	"github.com/example/cadastro-de-usuarios/internal/infrastructure/service"
	"github.com/example/cadastro-de-usuarios/internal/server"
)

func main() {
	db, err := infrastructure.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	deps := server.Dependencies{
		DB:           db,
		EmailSender:  service.NewEmailSender(),
		JWTValidator: service.NewJWTValidatorService(),
		UploadPath:   "./uploads",
	}

	e := server.SetupServer(deps)

	port := ":8080"
	log.Printf("Server listening on port %s\n", port)
	log.Fatal(e.Start(port))
}

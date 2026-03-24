package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
	"github.com/example/cadastro-de-usuarios/internal/infrastructure"
	"github.com/example/cadastro-de-usuarios/internal/infrastructure/repository"
	"github.com/example/cadastro-de-usuarios/internal/infrastructure/service"
	"github.com/example/cadastro-de-usuarios/internal/presentation"
	"github.com/example/cadastro-de-usuarios/internal/presentation/handler"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := infrastructure.Connect(databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	passwordRecoveryRepo := repository.NewPasswordRecoveryRepository(db)
	postRepo := repository.NewPostRepository(db)

	emailSender := service.NewEmailSender()
	jwtValidatorService := service.NewJWTValidatorService()

	registerUserUC := usecase.NewRegisterUserUseCase(userRepo)
	listUsersUC := usecase.NewListUsersUseCase(userRepo)
	deleteUserUC := usecase.NewDeleteUserUseCase(userRepo)
	updateUserProfileUC := usecase.NewUpdateUserProfileUseCase(userRepo)
	uploadProfilePictureUC := usecase.NewUploadProfilePictureUseCase(userRepo, "./uploads")

	requestPasswordRecoveryUC := usecase.NewRequestPasswordRecoveryUseCase(userRepo, passwordRecoveryRepo, emailSender)
	resetPasswordUC := usecase.NewResetPasswordUseCase(userRepo, passwordRecoveryRepo)
	createPostUC := usecase.NewCreatePostUseCase(postRepo)
	updatePostUC := usecase.NewUpdatePostUseCase(postRepo)
	validateTokenUC := usecase.NewValidateTokenUseCase(jwtValidatorService)

	userHandler := handler.NewUserHandler(registerUserUC, listUsersUC, updateUserProfileUC, deleteUserUC, uploadProfilePictureUC)
	passwordRecoveryHandler := handler.NewPasswordRecoveryHandler(requestPasswordRecoveryUC, resetPasswordUC)
	postHandler := handler.NewPostHandler(createPostUC, updatePostUC)

	e := echo.New()

	deps := presentation.RouterDependencies{
		UserHandler:             userHandler,
		PostHandler:             postHandler,
		PasswordRecoveryHandler: passwordRecoveryHandler,
	}
	presentation.NewRouter(e, deps, validateTokenUC)

	go func() {
		port := ":8080"
		log.Printf("Server listening on port %s\n", port)
		if err := e.Start(port); err != nil {
			log.Printf("Server stopped: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v\n", err)
	}

	log.Println("Server exited")
}

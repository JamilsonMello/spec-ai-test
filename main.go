package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/example/cadastro-de-usuarios/application/usecase"
	pkgdb "github.com/example/cadastro-de-usuarios/pkg/db"
	"github.com/example/cadastro-de-usuarios/infrastructure/repository"
	"github.com/example/cadastro-de-usuarios/infrastructure/service"
	"github.com/example/cadastro-de-usuarios/presentation/handler"
	"github.com/example/cadastro-de-usuarios/presentation/middleware"
)

func main() {
	db, err := pkgdb.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewPostgreSQLUserRepository(db)
	passwordRecoveryRepo := repository.NewPostgreSQLPasswordRecoveryRepository(db)
	postRepo := repository.NewPostgreSQLPostRepository(db)

	logger := log.New(os.Stdout, "", log.LstdFlags)

	emailCfg, err := service.LoadEmailServiceConfig()
	if err != nil {
		log.Fatalf("failed to load email service config: %v", err)
	}

	emailService, err := service.NewEmailService(emailCfg, logger)
	if err != nil {
		log.Fatalf("failed to create email service: %v", err)
	}
	emailService.StartWorkerPool()

	jwtValidatorService := service.NewJWTValidatorService()

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

	e := echo.New()

	e.POST("/usuarios", userHandler.RegisterUser)
	e.POST("/password-recovery", passwordRecoveryHandler.RequestPasswordRecovery)
	e.POST("/password-recovery/reset", passwordRecoveryHandler.ResetPassword)

	protected := e.Group("")
	protected.Use(middleware.AuthMiddleware(validateTokenUC))

	protected.GET("/usuarios/listar", userHandler.ListUsers)
	protected.DELETE("/usuarios/:id", userHandler.DeleteUser)
	protected.PUT("/usuarios/:id", userHandler.UpdateUserProfile)
	protected.POST("/posts", postHandler.CreatePost)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		emailService.StopWorkerPool()
		_ = e.Close()
	}()

	port := ":8080"
	log.Printf("Server listening on port %s\n", port)
	log.Fatal(e.Start(port))
}

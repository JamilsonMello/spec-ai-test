package main

import (
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
	"github.com/example/cadastro-de-usuarios/internal/infrastructure"
	"github.com/example/cadastro-de-usuarios/internal/infrastructure/repository"
	"github.com/example/cadastro-de-usuarios/internal/infrastructure/service"
	"github.com/example/cadastro-de-usuarios/internal/infrastructure/storage"
	"github.com/example/cadastro-de-usuarios/internal/presentation/handler"
	"github.com/example/cadastro-de-usuarios/internal/presentation/middleware"
)

func main() {
	db, err := infrastructure.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := infrastructure.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	passwordRecoveryRepo := repository.NewPasswordRecoveryRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	reactionRepo := repository.NewReactionRepository(db)
	communityRepo := repository.NewCommunityRepository(db)
	productRepo := repository.NewProductRepository(db)

	emailSender := service.NewEmailSender()
	bcryptHasher := service.NewBcryptHasher()
	txManager := infrastructure.NewSQLTransactionManager(db)
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

	userHandler := handler.NewUserHandler(registerUserUC, listUsersUC, updateUserProfileUC, deleteUserUC, uploadProfilePictureUC)
	passwordRecoveryHandler := handler.NewPasswordRecoveryHandler(requestPasswordRecoveryUC, resetPasswordUC)
	postHandler := handler.NewPostHandler(createPostUC, updatePostUC)
	commentHandler := handler.NewCommentHandler(createCommentUC, listCommentsUC)
	reactionHandler := handler.NewReactionHandler(toggleReactionUC)
	communityHandler := handler.NewCommunityHandler(createCommunityUC)
	productHandler := handler.NewProductHandler(createProductUC, updateProductUC, deleteProductUC)

	e := echo.New()

	e.Static("/uploads", "./uploads")

	e.POST("/usuarios", userHandler.RegisterUser)
	rateLimiter := middleware.RateLimitMiddleware(5, 1*time.Minute)
	e.POST("/password-recovery", passwordRecoveryHandler.RequestPasswordRecovery, rateLimiter)
	e.POST("/password-recovery/reset", passwordRecoveryHandler.ResetPassword)

	protected := e.Group("")
	protected.Use(middleware.AuthMiddleware(validateTokenUC))

	protected.GET("/usuarios/listar", userHandler.ListUsers)
	protected.DELETE("/usuarios/:id", userHandler.DeleteUser)
	protected.PUT("/usuarios/:id", userHandler.UpdateUserProfile)
	protected.POST("/usuarios/:id/foto-perfil", userHandler.UploadProfilePicture)
	protected.POST("/posts", postHandler.CreatePost)
	protected.PUT("/posts/:id", postHandler.UpdatePost)
	protected.POST("/posts/:id/comments", commentHandler.CreateComment)
	protected.POST("/comments/:id/reactions", reactionHandler.ToggleReaction)
	protected.POST("/comunidades", communityHandler.CreateCommunity)
	protected.POST("/products", productHandler.CreateProduct)
	protected.PUT("/products/:id", productHandler.UpdateProduct)
	protected.DELETE("/products/:id", productHandler.DeleteProduct)

	e.GET("/posts/:id/comments", commentHandler.ListComments)

	port := ":8080"
	log.Printf("Server listening on port %s\n", port)
	log.Fatal(e.Start(port))
}

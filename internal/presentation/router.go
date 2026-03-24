package presentation

import (
	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
	"github.com/example/cadastro-de-usuarios/internal/presentation/handler"
	"github.com/example/cadastro-de-usuarios/internal/presentation/middleware"
	"github.com/labstack/echo/v4"
)

type RouterDependencies struct {
	UserHandler             *handler.UserHandler
	PostHandler             *handler.PostHandler
	PasswordRecoveryHandler *handler.PasswordRecoveryHandler
}

func NewRouter(e *echo.Echo, deps RouterDependencies, validateTokenUC *usecase.ValidateTokenUseCase) {
	e.Static("/uploads", "./uploads")

	e.POST("/usuarios", deps.UserHandler.RegisterUser)
	e.POST("/password-recovery", deps.PasswordRecoveryHandler.RequestPasswordRecovery)
	e.POST("/password-recovery/reset", deps.PasswordRecoveryHandler.ResetPassword)

	protected := e.Group("")
	protected.Use(middleware.AuthMiddleware(validateTokenUC))

	protected.GET("/usuarios/listar", deps.UserHandler.ListUsers)
	protected.DELETE("/usuarios/:id", deps.UserHandler.DeleteUser)
	protected.PUT("/usuarios/:id", deps.UserHandler.UpdateUserProfile)
	protected.POST("/usuarios/:id/foto-perfil", deps.UserHandler.UploadProfilePicture)
	protected.POST("/posts", deps.PostHandler.CreatePost)
	protected.PUT("/posts/:id", deps.PostHandler.UpdatePost)
}

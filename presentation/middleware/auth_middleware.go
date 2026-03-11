package middleware

import (
	"net/http"
	"strings"

	"github.com/example/cadastro-de-usuarios/application/usecase"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(validateUC *usecase.ValidateTokenUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized access"})
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			resp, err := validateUC.Execute(usecase.ValidateTokenInput{TokenString: token})
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized access"})
			}

			c.Set("user_id", resp.Subject)

			return next(c)
		}
	}
}

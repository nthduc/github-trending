package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nthduc/github-trending/models"
	"github.com/nthduc/github-trending/security"
)

func JWTMiddleware() echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:     &models.JwtCustomClaims{},
		SigningKey: []byte(security.JWT_KEY),
	}

	return middleware.JWTWithConfig(config)
}

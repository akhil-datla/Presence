package server

import (
	"strings"

	"github.com/akhil-datla/Presence/internal/auth"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware validates the Authorization header and sets user_id in context.
func JWTMiddleware(jwt *auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return echo.NewHTTPError(401, map[string]string{"error": "missing authorization header"})
			}

			token := strings.TrimPrefix(header, "Bearer ")
			if token == header {
				return echo.NewHTTPError(401, map[string]string{"error": "invalid authorization format"})
			}

			claims, err := jwt.ValidateToken(token)
			if err != nil {
				return echo.NewHTTPError(401, map[string]string{"error": "invalid or expired token"})
			}

			c.Set("user_id", claims.UserID)
			return next(c)
		}
	}
}

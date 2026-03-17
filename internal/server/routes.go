package server

import (
	"net/http"

	"github.com/akhil-datla/Presence/internal/auth"
	"github.com/akhil-datla/Presence/internal/handler"
	"github.com/akhil-datla/Presence/internal/store"
	"github.com/labstack/echo/v4"
)

func registerRoutes(e *echo.Echo, s *store.Store, jwt *auth.JWTService) {
	authH := handler.NewAuthHandler(s, jwt)
	userH := handler.NewUserHandler(s)
	sessH := handler.NewSessionHandler(s)
	attH := handler.NewAttendanceHandler(s)

	api := e.Group("/api/v1")

	// Health check
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Public routes
	api.POST("/auth/register", authH.Register)
	api.POST("/auth/login", authH.Login)

	// Protected routes
	protected := api.Group("", JWTMiddleware(jwt))

	// User
	protected.GET("/users/me", userH.GetProfile)
	protected.PUT("/users/me", userH.UpdateProfile)
	protected.DELETE("/users/me", userH.DeleteProfile)

	// Sessions
	protected.POST("/sessions", sessH.Create)
	protected.GET("/sessions", sessH.List)
	protected.GET("/sessions/:id", sessH.Get)
	protected.PUT("/sessions/:id", sessH.Update)
	protected.DELETE("/sessions/:id", sessH.Delete)

	// Attendance
	protected.POST("/sessions/:id/checkin", attH.CheckIn)
	protected.POST("/sessions/:id/checkout", attH.CheckOut)
	protected.GET("/sessions/:id/attendance", attH.List)
	protected.GET("/sessions/:id/attendance/filter", attH.Filter)
	protected.DELETE("/sessions/:id/attendance", attH.Clear)
	protected.GET("/sessions/:id/export/csv", attH.ExportCSV)
}

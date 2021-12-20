package apiserver

import (
	"main/internal/attendance"
	"main/internal/organizers"
	"main/internal/participants"
	"main/internal/sessions"
	"net/http"

	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var e *echo.Echo

//Start the server
func Start(portNum string) {
	e = echo.New()
	e.HideBanner = true
	organizers.New()
	participants.New()
	sessions.New()
	attendance.New()
	InitRoutes()
	e.Logger.SetLevel(log.INFO)
	e.Logger.Fatal(e.Start(portNum))
}

//InitRoutes initializes the routes for the server
func InitRoutes() {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}

	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))
	e.POST("/organizer/add", AddOrganizer)
	e.POST("/organizer/view", ViewOrganizer)
	e.POST("/organizer/login", AuthenticateOrganizer)
	e.POST("/organizer/update", UpdateOrganizer)
	e.POST("/organizer/delete", DeleteOrganizer)
	e.POST("/participant/add", AddParticipant)
	e.POST("/participant/view", ViewParticipant)
	e.POST("/participant/login", AuthenticateParticipant)
	e.POST("/particpant/update", UpdateOrganizer)
	e.POST("/participant/delete", DeleteOrganizer)
	e.POST("/session/create", CreateSession)
	e.POST("/session/view", GetSessions)
	e.POST("/session/update", UpdateSession)
	e.POST("/session/delete", DeleteSession)
	e.POST("/session/attendance", GetAttendance)
	e.POST("/session/attendance/filter", FilterAttendance)
	e.POST("/attendance/in", CheckIn)
	e.POST("/attendance/out", CheckOut)
	e.POST("/attendance/clear", ClearAttendance)
}

//Close the server
func Close() error {
	err := e.Close()
	return err
}

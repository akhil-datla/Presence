package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type successResponse struct {
	Data interface{} `json:"data"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func ok(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, successResponse{Data: data})
}

func created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, successResponse{Data: data})
}

func badRequest(c echo.Context, msg string) error {
	return c.JSON(http.StatusBadRequest, errorResponse{Error: msg})
}

func unauthorized(c echo.Context, msg string) error {
	return c.JSON(http.StatusUnauthorized, errorResponse{Error: msg})
}

func notFound(c echo.Context, msg string) error {
	return c.JSON(http.StatusNotFound, errorResponse{Error: msg})
}

func conflict(c echo.Context, msg string) error {
	return c.JSON(http.StatusConflict, errorResponse{Error: msg})
}

func serverError(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, errorResponse{Error: "internal server error"})
}

package handler

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/akhil-datla/Presence/internal/auth"
	"github.com/akhil-datla/Presence/internal/model"
	"github.com/akhil-datla/Presence/internal/store"
	"github.com/labstack/echo/v4"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	store *store.Store
	jwt   *auth.JWTService
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(s *store.Store, jwt *auth.JWTService) *AuthHandler {
	return &AuthHandler{store: s, jwt: jwt}
}

type tokenResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// Register creates a new user account.
func (h *AuthHandler) Register(c echo.Context) error {
	var req model.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if err := req.Validate(); err != nil {
		return badRequest(c, err.Error())
	}

	hashed, err := auth.HashPassword(req.Password)
	if err != nil {
		return serverError(c)
	}

	user, err := h.store.CreateUser(req.FirstName, req.LastName, req.Email, hashed)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return conflict(c, "email already registered")
		}
		return serverError(c)
	}

	token, err := h.jwt.GenerateToken(user.ID)
	if err != nil {
		return serverError(c)
	}

	user.Password = ""
	return created(c, tokenResponse{Token: token, User: user})
}

// Login authenticates a user and returns a JWT.
func (h *AuthHandler) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if err := req.Validate(); err != nil {
		return badRequest(c, err.Error())
	}

	user, err := h.store.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return unauthorized(c, "invalid credentials")
		}
		return serverError(c)
	}

	if !auth.CheckPassword(user.Password, req.Password) {
		return unauthorized(c, "invalid credentials")
	}

	token, err := h.jwt.GenerateToken(user.ID)
	if err != nil {
		return serverError(c)
	}

	user.Password = ""
	return ok(c, tokenResponse{Token: token, User: user})
}

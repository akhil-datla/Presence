package handler

import (
	"github.com/akhil-datla/Presence/internal/auth"
	"github.com/akhil-datla/Presence/internal/model"
	"github.com/akhil-datla/Presence/internal/store"
	"github.com/labstack/echo/v4"
)

// UserHandler handles user profile endpoints.
type UserHandler struct {
	store *store.Store
}

// NewUserHandler creates a new user handler.
func NewUserHandler(s *store.Store) *UserHandler {
	return &UserHandler{store: s}
}

// GetProfile returns the current user's profile.
func (h *UserHandler) GetProfile(c echo.Context) error {
	user, err := h.store.GetUserByID(userID(c))
	if err != nil {
		return notFound(c, "user not found")
	}
	user.Password = ""
	return ok(c, user)
}

// UpdateProfile updates the current user's profile.
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	var req model.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if err := req.Validate(); err != nil {
		return badRequest(c, err.Error())
	}

	var hashed *string
	if req.Password != nil {
		h, err := auth.HashPassword(*req.Password)
		if err != nil {
			return serverError(c)
		}
		hashed = &h
	}

	user, err := h.store.UpdateUser(userID(c), &req, hashed)
	if err != nil {
		return serverError(c)
	}
	user.Password = ""
	return ok(c, user)
}

// DeleteProfile deletes the current user's account.
func (h *UserHandler) DeleteProfile(c echo.Context) error {
	if err := h.store.DeleteUser(userID(c)); err != nil {
		return serverError(c)
	}
	return ok(c, map[string]string{"message": "account deleted"})
}

// userID extracts the authenticated user's ID from the context.
func userID(c echo.Context) string {
	return c.Get("user_id").(string)
}

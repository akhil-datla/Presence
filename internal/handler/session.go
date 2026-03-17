package handler

import (
	"database/sql"
	"errors"

	"github.com/akhil-datla/Presence/internal/model"
	"github.com/akhil-datla/Presence/internal/store"
	"github.com/labstack/echo/v4"
)

// SessionHandler handles session CRUD endpoints.
type SessionHandler struct {
	store *store.Store
}

// NewSessionHandler creates a new session handler.
func NewSessionHandler(s *store.Store) *SessionHandler {
	return &SessionHandler{store: s}
}

// Create creates a new session.
func (h *SessionHandler) Create(c echo.Context) error {
	var req model.CreateSessionRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if err := req.Validate(); err != nil {
		return badRequest(c, err.Error())
	}

	sess, err := h.store.CreateSession(userID(c), req.Name)
	if err != nil {
		return serverError(c)
	}
	return created(c, sess)
}

// List returns all sessions for the authenticated user.
func (h *SessionHandler) List(c echo.Context) error {
	sessions, err := h.store.ListSessions(userID(c))
	if err != nil {
		return serverError(c)
	}
	if sessions == nil {
		sessions = []model.Session{}
	}
	return ok(c, sessions)
}

// Get returns a single session by ID.
func (h *SessionHandler) Get(c echo.Context) error {
	sess, err := h.store.GetSession(c.Param("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return notFound(c, "session not found")
		}
		return serverError(c)
	}
	return ok(c, sess)
}

// Update updates a session.
func (h *SessionHandler) Update(c echo.Context) error {
	var req model.UpdateSessionRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if err := req.Validate(); err != nil {
		return badRequest(c, err.Error())
	}

	sess, err := h.store.UpdateSession(c.Param("id"), userID(c), *req.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return notFound(c, "session not found or not owned by you")
		}
		return serverError(c)
	}
	return ok(c, sess)
}

// Delete removes a session.
func (h *SessionHandler) Delete(c echo.Context) error {
	err := h.store.DeleteSession(c.Param("id"), userID(c))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return notFound(c, "session not found or not owned by you")
		}
		return serverError(c)
	}
	return ok(c, map[string]string{"message": "session deleted"})
}

package model

import (
	"errors"
	"time"
)

// Session represents an attendance session created by an organizer.
type Session struct {
	ID          string    `json:"id"`
	OrganizerID string    `json:"organizer_id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateSessionRequest is the payload for creating a session.
type CreateSessionRequest struct {
	Name string `json:"name"`
}

// Validate checks the create session request.
func (r *CreateSessionRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// UpdateSessionRequest is the payload for updating a session.
type UpdateSessionRequest struct {
	Name *string `json:"name,omitempty"`
}

// Validate checks the update session request.
func (r *UpdateSessionRequest) Validate() error {
	if r.Name != nil && *r.Name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

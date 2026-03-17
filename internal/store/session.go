package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/akhil-datla/Presence/internal/model"
)

// CreateSession inserts a new session.
func (s *Store) CreateSession(organizerID, name string) (*model.Session, error) {
	id := newID()
	now := time.Now().UTC()

	_, err := s.db.Exec(
		`INSERT INTO sessions (id, organizer_id, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		id, organizerID, name, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &model.Session{
		ID: id, OrganizerID: organizerID, Name: name,
		CreatedAt: now, UpdatedAt: now,
	}, nil
}

// GetSession returns a session by ID.
func (s *Store) GetSession(id string) (*model.Session, error) {
	var sess model.Session
	err := s.db.QueryRow(
		`SELECT id, organizer_id, name, created_at, updated_at FROM sessions WHERE id = ?`, id,
	).Scan(&sess.ID, &sess.OrganizerID, &sess.Name, &sess.CreatedAt, &sess.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// ListSessions returns all sessions for an organizer.
func (s *Store) ListSessions(organizerID string) ([]model.Session, error) {
	rows, err := s.db.Query(
		`SELECT id, organizer_id, name, created_at, updated_at FROM sessions WHERE organizer_id = ? ORDER BY created_at DESC`, organizerID,
	)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	defer rows.Close()

	var sessions []model.Session
	for rows.Next() {
		var sess model.Session
		if err := rows.Scan(&sess.ID, &sess.OrganizerID, &sess.Name, &sess.CreatedAt, &sess.UpdatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, sess)
	}
	return sessions, rows.Err()
}

// UpdateSession updates a session's name.
func (s *Store) UpdateSession(id, organizerID, name string) (*model.Session, error) {
	now := time.Now().UTC()
	res, err := s.db.Exec(
		`UPDATE sessions SET name = ?, updated_at = ? WHERE id = ? AND organizer_id = ?`,
		name, now, id, organizerID,
	)
	if err != nil {
		return nil, fmt.Errorf("update session: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, sql.ErrNoRows
	}
	return s.GetSession(id)
}

// DeleteSession removes a session.
func (s *Store) DeleteSession(id, organizerID string) error {
	res, err := s.db.Exec(`DELETE FROM sessions WHERE id = ? AND organizer_id = ?`, id, organizerID)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

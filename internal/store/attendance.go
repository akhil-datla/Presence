package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/akhil-datla/Presence/internal/model"
)

// CheckIn records a participant checking in to a session.
func (s *Store) CheckIn(sessionID, participantID, participantName string) (*model.Attendance, error) {
	id := newID()
	now := time.Now().UTC()

	_, err := s.db.Exec(
		`INSERT INTO attendance (id, session_id, participant_id, participant_name, time_in) VALUES (?, ?, ?, ?, ?)`,
		id, sessionID, participantID, participantName, now,
	)
	if err != nil {
		return nil, fmt.Errorf("check in: %w", err)
	}

	return &model.Attendance{
		ID: id, SessionID: sessionID,
		ParticipantID: participantID, ParticipantName: participantName,
		TimeIn: now,
	}, nil
}

// CheckOut records a participant checking out. Updates the most recent open check-in.
func (s *Store) CheckOut(sessionID, participantID string) (*model.Attendance, error) {
	now := time.Now().UTC()

	// Find the most recent check-in without a check-out
	var att model.Attendance
	var timeOut sql.NullTime
	err := s.db.QueryRow(
		`SELECT id, session_id, participant_id, participant_name, time_in, time_out
		 FROM attendance
		 WHERE session_id = ? AND participant_id = ? AND time_out IS NULL
		 ORDER BY time_in DESC LIMIT 1`,
		sessionID, participantID,
	).Scan(&att.ID, &att.SessionID, &att.ParticipantID, &att.ParticipantName, &att.TimeIn, &timeOut)
	if err != nil {
		return nil, fmt.Errorf("no open check-in found: %w", err)
	}

	_, err = s.db.Exec(`UPDATE attendance SET time_out = ? WHERE id = ?`, now, att.ID)
	if err != nil {
		return nil, fmt.Errorf("check out: %w", err)
	}

	att.TimeOut = &now
	return &att, nil
}

// GetAttendance returns all attendance records for a session.
func (s *Store) GetAttendance(sessionID string) ([]model.Attendance, error) {
	rows, err := s.db.Query(
		`SELECT id, session_id, participant_id, participant_name, time_in, time_out
		 FROM attendance WHERE session_id = ? ORDER BY time_in DESC`, sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("get attendance: %w", err)
	}
	defer func() { _ = rows.Close() }()

	return scanAttendanceRows(rows)
}

// FilterAttendance returns attendance records filtered by time.
func (s *Store) FilterAttendance(sessionID string, t time.Time, mode string) ([]model.Attendance, error) {
	var query string
	switch mode {
	case "before":
		query = `SELECT id, session_id, participant_id, participant_name, time_in, time_out
		         FROM attendance WHERE session_id = ? AND time_in < ? ORDER BY time_in DESC`
	case "after":
		query = `SELECT id, session_id, participant_id, participant_name, time_in, time_out
		         FROM attendance WHERE session_id = ? AND time_in > ? ORDER BY time_in DESC`
	default:
		return nil, fmt.Errorf("invalid filter mode: %q (use \"before\" or \"after\")", mode)
	}

	rows, err := s.db.Query(query, sessionID, t)
	if err != nil {
		return nil, fmt.Errorf("filter attendance: %w", err)
	}
	defer func() { _ = rows.Close() }()

	return scanAttendanceRows(rows)
}

// ClearAttendance deletes all attendance records for a session.
func (s *Store) ClearAttendance(sessionID, organizerID string) error {
	// Verify the session belongs to the organizer
	var id string
	err := s.db.QueryRow(`SELECT id FROM sessions WHERE id = ? AND organizer_id = ?`, sessionID, organizerID).Scan(&id)
	if err != nil {
		return fmt.Errorf("session not found or not owned: %w", err)
	}

	_, err = s.db.Exec(`DELETE FROM attendance WHERE session_id = ?`, sessionID)
	return err
}

func scanAttendanceRows(rows *sql.Rows) ([]model.Attendance, error) {
	var list []model.Attendance
	for rows.Next() {
		var att model.Attendance
		var timeOut sql.NullTime
		if err := rows.Scan(&att.ID, &att.SessionID, &att.ParticipantID, &att.ParticipantName, &att.TimeIn, &timeOut); err != nil {
			return nil, err
		}
		if timeOut.Valid {
			att.TimeOut = &timeOut.Time
		}
		list = append(list, att)
	}
	return list, rows.Err()
}

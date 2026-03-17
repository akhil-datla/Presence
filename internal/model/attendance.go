package model

import "time"

// Attendance represents a check-in/check-out record for a session participant.
type Attendance struct {
	ID              string     `json:"id"`
	SessionID       string     `json:"session_id"`
	ParticipantID   string     `json:"participant_id"`
	ParticipantName string     `json:"participant_name"`
	TimeIn          time.Time  `json:"time_in"`
	TimeOut         *time.Time `json:"time_out,omitempty"`
}

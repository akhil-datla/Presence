package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/akhil-datla/Presence/internal/model"
)

// CreateUser inserts a new user.
func (s *Store) CreateUser(firstName, lastName, email, hashedPassword string) (*model.User, error) {
	id := newID()
	now := time.Now().UTC()

	_, err := s.db.Exec(
		`INSERT INTO users (id, first_name, last_name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, firstName, lastName, email, hashedPassword, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &model.User{
		ID: id, FirstName: firstName, LastName: lastName,
		Email: email, Password: hashedPassword,
		CreatedAt: now, UpdatedAt: now,
	}, nil
}

// GetUserByID returns a user by ID.
func (s *Store) GetUserByID(id string) (*model.User, error) {
	return s.scanUser(s.db.QueryRow(`SELECT id, first_name, last_name, email, password, created_at, updated_at FROM users WHERE id = ?`, id))
}

// GetUserByEmail returns a user by email.
func (s *Store) GetUserByEmail(email string) (*model.User, error) {
	return s.scanUser(s.db.QueryRow(`SELECT id, first_name, last_name, email, password, created_at, updated_at FROM users WHERE email = ?`, email))
}

// UpdateUser updates the specified fields of a user.
func (s *Store) UpdateUser(id string, req *model.UpdateUserRequest, hashedPassword *string) (*model.User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if hashedPassword != nil {
		user.Password = *hashedPassword
	}
	user.UpdatedAt = time.Now().UTC()

	_, err = s.db.Exec(
		`UPDATE users SET first_name = ?, last_name = ?, email = ?, password = ?, updated_at = ? WHERE id = ?`,
		user.FirstName, user.LastName, user.Email, user.Password, user.UpdatedAt, id,
	)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

// DeleteUser removes a user by ID.
func (s *Store) DeleteUser(id string) error {
	res, err := s.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Store) scanUser(row *sql.Row) (*model.User, error) {
	var u model.User
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate ID: " + err.Error())
	}
	return hex.EncodeToString(b)
}

package store

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/akhil-datla/Presence/internal/model"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func createTestUser(t *testing.T, s *Store) *model.User {
	t.Helper()
	u, err := s.CreateUser("John", "Doe", "john@example.com", "hashedpw123")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return u
}

// ---------- User Tests ----------

func TestCreateUser(t *testing.T) {
	s := newTestStore(t)

	t.Run("success", func(t *testing.T) {
		u, err := s.CreateUser("Alice", "Smith", "alice@test.com", "hashed")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if u.ID == "" {
			t.Fatal("expected non-empty ID")
		}
		if u.FirstName != "Alice" {
			t.Fatalf("expected FirstName=Alice, got %s", u.FirstName)
		}
		if u.LastName != "Smith" {
			t.Fatalf("expected LastName=Smith, got %s", u.LastName)
		}
		if u.Email != "alice@test.com" {
			t.Fatalf("expected Email=alice@test.com, got %s", u.Email)
		}
		if u.Password != "hashed" {
			t.Fatalf("expected Password=hashed, got %s", u.Password)
		}
		if u.CreatedAt.IsZero() {
			t.Fatal("expected non-zero CreatedAt")
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		_, err := s.CreateUser("Bob", "Jones", "alice@test.com", "hashed")
		if err == nil {
			t.Fatal("expected error for duplicate email")
		}
	})
}

func TestGetUserByID(t *testing.T) {
	s := newTestStore(t)
	u := createTestUser(t, s)

	t.Run("existing user", func(t *testing.T) {
		found, err := s.GetUserByID(u.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if found.ID != u.ID {
			t.Fatalf("expected ID=%s, got %s", u.ID, found.ID)
		}
		if found.Email != u.Email {
			t.Fatalf("expected Email=%s, got %s", u.Email, found.Email)
		}
	})

	t.Run("non-existent user", func(t *testing.T) {
		_, err := s.GetUserByID("nonexistent")
		if err == nil {
			t.Fatal("expected error for non-existent user")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("expected sql.ErrNoRows, got %v", err)
		}
	})
}

func TestGetUserByEmail(t *testing.T) {
	s := newTestStore(t)
	u := createTestUser(t, s)

	t.Run("existing email", func(t *testing.T) {
		found, err := s.GetUserByEmail(u.Email)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if found.ID != u.ID {
			t.Fatalf("expected ID=%s, got %s", u.ID, found.ID)
		}
	})

	t.Run("non-existent email", func(t *testing.T) {
		_, err := s.GetUserByEmail("nobody@test.com")
		if err == nil {
			t.Fatal("expected error for non-existent email")
		}
	})
}

func TestUpdateUser(t *testing.T) {
	s := newTestStore(t)
	u := createTestUser(t, s)

	strPtr := func(s string) *string { return &s }

	t.Run("update first name", func(t *testing.T) {
		req := &model.UpdateUserRequest{FirstName: strPtr("Jane")}
		updated, err := s.UpdateUser(u.ID, req, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.FirstName != "Jane" {
			t.Fatalf("expected FirstName=Jane, got %s", updated.FirstName)
		}
		if updated.LastName != u.LastName {
			t.Fatalf("expected LastName unchanged, got %s", updated.LastName)
		}
	})

	t.Run("update password", func(t *testing.T) {
		req := &model.UpdateUserRequest{}
		newHash := "newhash"
		updated, err := s.UpdateUser(u.ID, req, &newHash)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.Password != "newhash" {
			t.Fatalf("expected Password=newhash, got %s", updated.Password)
		}
	})

	t.Run("update multiple fields", func(t *testing.T) {
		req := &model.UpdateUserRequest{
			LastName: strPtr("Updated"),
			Email:    strPtr("updated@test.com"),
		}
		updated, err := s.UpdateUser(u.ID, req, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.LastName != "Updated" {
			t.Fatalf("expected LastName=Updated, got %s", updated.LastName)
		}
		if updated.Email != "updated@test.com" {
			t.Fatalf("expected Email=updated@test.com, got %s", updated.Email)
		}
	})

	t.Run("non-existent user", func(t *testing.T) {
		req := &model.UpdateUserRequest{FirstName: strPtr("X")}
		_, err := s.UpdateUser("nonexistent", req, nil)
		if err == nil {
			t.Fatal("expected error for non-existent user")
		}
	})
}

func TestDeleteUser(t *testing.T) {
	s := newTestStore(t)

	t.Run("delete existing user", func(t *testing.T) {
		u, err := s.CreateUser("Del", "User", "del@test.com", "hashed")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := s.DeleteUser(u.ID); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		_, err = s.GetUserByID(u.ID)
		if err == nil {
			t.Fatal("expected error after deletion")
		}
	})

	t.Run("delete non-existent user", func(t *testing.T) {
		err := s.DeleteUser("nonexistent")
		if err == nil {
			t.Fatal("expected error for non-existent user")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("expected sql.ErrNoRows, got %v", err)
		}
	})
}

// ---------- Session Tests ----------

func TestCreateSession(t *testing.T) {
	s := newTestStore(t)
	u := createTestUser(t, s)

	t.Run("success", func(t *testing.T) {
		sess, err := s.CreateSession(u.ID, "Test Session")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sess.ID == "" {
			t.Fatal("expected non-empty ID")
		}
		if sess.OrganizerID != u.ID {
			t.Fatalf("expected OrganizerID=%s, got %s", u.ID, sess.OrganizerID)
		}
		if sess.Name != "Test Session" {
			t.Fatalf("expected Name=Test Session, got %s", sess.Name)
		}
		if sess.CreatedAt.IsZero() {
			t.Fatal("expected non-zero CreatedAt")
		}
	})
}

func TestGetSession(t *testing.T) {
	s := newTestStore(t)
	u := createTestUser(t, s)

	sess, err := s.CreateSession(u.ID, "My Session")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("existing session", func(t *testing.T) {
		found, err := s.GetSession(sess.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if found.Name != "My Session" {
			t.Fatalf("expected Name=My Session, got %s", found.Name)
		}
	})

	t.Run("non-existent session", func(t *testing.T) {
		_, err := s.GetSession("nonexistent")
		if err == nil {
			t.Fatal("expected error for non-existent session")
		}
	})
}

func TestListSessions(t *testing.T) {
	s := newTestStore(t)
	u := createTestUser(t, s)

	t.Run("empty list", func(t *testing.T) {
		sessions, err := s.ListSessions(u.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sessions) != 0 {
			t.Fatalf("expected 0 sessions, got %d", len(sessions))
		}
	})

	_, _ = s.CreateSession(u.ID, "Session 1")
	_, _ = s.CreateSession(u.ID, "Session 2")

	t.Run("multiple sessions", func(t *testing.T) {
		sessions, err := s.ListSessions(u.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sessions) != 2 {
			t.Fatalf("expected 2 sessions, got %d", len(sessions))
		}
	})

	t.Run("different organizer sees none", func(t *testing.T) {
		other, _ := s.CreateUser("Other", "User", "other@test.com", "hashed")
		sessions, err := s.ListSessions(other.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sessions) != 0 {
			t.Fatalf("expected 0 sessions for other user, got %d", len(sessions))
		}
	})
}

func TestUpdateSession(t *testing.T) {
	s := newTestStore(t)
	u := createTestUser(t, s)
	sess, _ := s.CreateSession(u.ID, "Original")

	t.Run("success", func(t *testing.T) {
		updated, err := s.UpdateSession(sess.ID, u.ID, "Renamed")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.Name != "Renamed" {
			t.Fatalf("expected Name=Renamed, got %s", updated.Name)
		}
	})

	t.Run("wrong organizer", func(t *testing.T) {
		_, err := s.UpdateSession(sess.ID, "wronguser", "Hack")
		if err == nil {
			t.Fatal("expected error for wrong organizer")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("expected sql.ErrNoRows, got %v", err)
		}
	})

	t.Run("non-existent session", func(t *testing.T) {
		_, err := s.UpdateSession("nonexistent", u.ID, "X")
		if err == nil {
			t.Fatal("expected error for non-existent session")
		}
	})
}

func TestDeleteSession(t *testing.T) {
	s := newTestStore(t)
	u := createTestUser(t, s)

	t.Run("success", func(t *testing.T) {
		sess, _ := s.CreateSession(u.ID, "ToDelete")
		if err := s.DeleteSession(sess.ID, u.ID); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		_, err := s.GetSession(sess.ID)
		if err == nil {
			t.Fatal("expected error after deletion")
		}
	})

	t.Run("wrong organizer", func(t *testing.T) {
		sess, _ := s.CreateSession(u.ID, "Protected")
		err := s.DeleteSession(sess.ID, "wronguser")
		if err == nil {
			t.Fatal("expected error for wrong organizer")
		}
	})

	t.Run("non-existent session", func(t *testing.T) {
		err := s.DeleteSession("nonexistent", u.ID)
		if err == nil {
			t.Fatal("expected error for non-existent session")
		}
	})
}

// ---------- Attendance Tests ----------

func TestCheckIn(t *testing.T) {
	s := newTestStore(t)
	u := createTestUser(t, s)
	sess, _ := s.CreateSession(u.ID, "Attendance Session")

	t.Run("success", func(t *testing.T) {
		att, err := s.CheckIn(sess.ID, u.ID, "John Doe")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if att.ID == "" {
			t.Fatal("expected non-empty ID")
		}
		if att.SessionID != sess.ID {
			t.Fatalf("expected SessionID=%s, got %s", sess.ID, att.SessionID)
		}
		if att.ParticipantID != u.ID {
			t.Fatalf("expected ParticipantID=%s, got %s", u.ID, att.ParticipantID)
		}
		if att.ParticipantName != "John Doe" {
			t.Fatalf("expected ParticipantName=John Doe, got %s", att.ParticipantName)
		}
		if att.TimeIn.IsZero() {
			t.Fatal("expected non-zero TimeIn")
		}
		if att.TimeOut != nil {
			t.Fatal("expected nil TimeOut on check-in")
		}
	})
}

func TestCheckOut(t *testing.T) {
	s := newTestStore(t)
	u := createTestUser(t, s)
	sess, _ := s.CreateSession(u.ID, "Checkout Session")

	t.Run("success after check-in", func(t *testing.T) {
		_, err := s.CheckIn(sess.ID, u.ID, "John Doe")
		if err != nil {
			t.Fatalf("check-in error: %v", err)
		}

		att, err := s.CheckOut(sess.ID, u.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if att.TimeOut == nil {
			t.Fatal("expected non-nil TimeOut after check-out")
		}
	})

	t.Run("no open check-in", func(t *testing.T) {
		// The previous check-in was already checked out
		_, err := s.CheckOut(sess.ID, u.ID)
		if err == nil {
			t.Fatal("expected error when no open check-in exists")
		}
	})
}

func TestGetAttendance(t *testing.T) {
	s := newTestStore(t)
	u := createTestUser(t, s)
	sess, _ := s.CreateSession(u.ID, "Att Session")

	t.Run("empty attendance", func(t *testing.T) {
		records, err := s.GetAttendance(sess.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(records) != 0 {
			t.Fatalf("expected 0 records, got %d", len(records))
		}
	})

	_, _ = s.CheckIn(sess.ID, u.ID, "John Doe")

	t.Run("one record", func(t *testing.T) {
		records, err := s.GetAttendance(sess.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(records) != 1 {
			t.Fatalf("expected 1 record, got %d", len(records))
		}
		if records[0].ParticipantName != "John Doe" {
			t.Fatalf("expected ParticipantName=John Doe, got %s", records[0].ParticipantName)
		}
	})
}

func TestFilterAttendance(t *testing.T) {
	s := newTestStore(t)
	u := createTestUser(t, s)
	sess, _ := s.CreateSession(u.ID, "Filter Session")

	// Create a check-in (time is set to now by the store)
	_, _ = s.CheckIn(sess.ID, u.ID, "John Doe")

	t.Run("filter after past time returns record", func(t *testing.T) {
		pastTime := time.Now().Add(-1 * time.Hour)
		records, err := s.FilterAttendance(sess.ID, pastTime, "after")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(records) != 1 {
			t.Fatalf("expected 1 record, got %d", len(records))
		}
	})

	t.Run("filter before past time returns nothing", func(t *testing.T) {
		pastTime := time.Now().Add(-1 * time.Hour)
		records, err := s.FilterAttendance(sess.ID, pastTime, "before")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(records) != 0 {
			t.Fatalf("expected 0 records, got %d", len(records))
		}
	})

	t.Run("filter before future time returns record", func(t *testing.T) {
		futureTime := time.Now().UTC().Add(1 * time.Hour)
		records, err := s.FilterAttendance(sess.ID, futureTime, "before")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(records) != 1 {
			t.Fatalf("expected 1 record, got %d", len(records))
		}
	})

	t.Run("invalid mode", func(t *testing.T) {
		_, err := s.FilterAttendance(sess.ID, time.Now(), "invalid")
		if err == nil {
			t.Fatal("expected error for invalid mode")
		}
	})
}

func TestClearAttendance(t *testing.T) {
	s := newTestStore(t)
	u := createTestUser(t, s)
	sess, _ := s.CreateSession(u.ID, "Clear Session")

	_, _ = s.CheckIn(sess.ID, u.ID, "John Doe")

	t.Run("success", func(t *testing.T) {
		err := s.ClearAttendance(sess.ID, u.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		records, err := s.GetAttendance(sess.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(records) != 0 {
			t.Fatalf("expected 0 records after clear, got %d", len(records))
		}
	})

	t.Run("wrong organizer", func(t *testing.T) {
		err := s.ClearAttendance(sess.ID, "wronguser")
		if err == nil {
			t.Fatal("expected error for wrong organizer")
		}
	})

	t.Run("non-existent session", func(t *testing.T) {
		err := s.ClearAttendance("nonexistent", u.ID)
		if err == nil {
			t.Fatal("expected error for non-existent session")
		}
	})
}

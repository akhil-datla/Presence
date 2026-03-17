package auth

import (
	"testing"
	"time"
)

func TestGenerateAndValidateToken(t *testing.T) {
	svc := NewJWTService("test-secret-key")

	t.Run("generate and validate token", func(t *testing.T) {
		token, err := svc.GenerateToken("user-123")
		if err != nil {
			t.Fatalf("unexpected error generating token: %v", err)
		}
		if token == "" {
			t.Fatal("expected non-empty token")
		}

		claims, err := svc.ValidateToken(token)
		if err != nil {
			t.Fatalf("unexpected error validating token: %v", err)
		}
		if claims.UserID != "user-123" {
			t.Fatalf("expected UserID=user-123, got %s", claims.UserID)
		}
	})

	t.Run("different users get different tokens", func(t *testing.T) {
		t1, _ := svc.GenerateToken("user-1")
		t2, _ := svc.GenerateToken("user-2")
		if t1 == t2 {
			t.Fatal("expected different tokens for different users")
		}
	})
}

func TestExpiredToken(t *testing.T) {
	// Create a service with a very short expiry
	svc := &JWTService{
		secret: []byte("test-secret"),
		expiry: -1 * time.Second, // already expired
	}

	token, err := svc.GenerateToken("user-123")
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	_, err = svc.ValidateToken(token)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestInvalidToken(t *testing.T) {
	svc := NewJWTService("test-secret-key")

	t.Run("garbage string", func(t *testing.T) {
		_, err := svc.ValidateToken("not-a-valid-token")
		if err == nil {
			t.Fatal("expected error for invalid token")
		}
	})

	t.Run("empty string", func(t *testing.T) {
		_, err := svc.ValidateToken("")
		if err == nil {
			t.Fatal("expected error for empty token")
		}
	})

	t.Run("wrong secret", func(t *testing.T) {
		otherSvc := NewJWTService("different-secret")
		token, _ := otherSvc.GenerateToken("user-123")

		_, err := svc.ValidateToken(token)
		if err == nil {
			t.Fatal("expected error for token signed with different secret")
		}
	})
}

func TestHashPassword(t *testing.T) {
	t.Run("hash and check", func(t *testing.T) {
		password := "mysecretpassword"
		hashed, err := HashPassword(password)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if hashed == "" {
			t.Fatal("expected non-empty hash")
		}
		if hashed == password {
			t.Fatal("hash should not equal plaintext password")
		}

		if !CheckPassword(hashed, password) {
			t.Fatal("CheckPassword should return true for correct password")
		}
	})

	t.Run("wrong password fails", func(t *testing.T) {
		hashed, _ := HashPassword("correct-password")
		if CheckPassword(hashed, "wrong-password") {
			t.Fatal("CheckPassword should return false for wrong password")
		}
	})

	t.Run("different hashes for same password", func(t *testing.T) {
		h1, _ := HashPassword("same-password")
		h2, _ := HashPassword("same-password")
		if h1 == h2 {
			t.Fatal("bcrypt should produce different hashes due to salt")
		}
		// Both should still validate
		if !CheckPassword(h1, "same-password") {
			t.Fatal("first hash should validate")
		}
		if !CheckPassword(h2, "same-password") {
			t.Fatal("second hash should validate")
		}
	})
}

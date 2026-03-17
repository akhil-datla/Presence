package model

import (
	"testing"
)

func TestCreateUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateUserRequest
		wantErr string
	}{
		{
			name:    "empty first name",
			req:     CreateUserRequest{FirstName: "", LastName: "Doe", Email: "a@b.com", Password: "12345678"},
			wantErr: "first_name is required",
		},
		{
			name:    "empty last name",
			req:     CreateUserRequest{FirstName: "John", LastName: "", Email: "a@b.com", Password: "12345678"},
			wantErr: "last_name is required",
		},
		{
			name:    "empty email",
			req:     CreateUserRequest{FirstName: "John", LastName: "Doe", Email: "", Password: "12345678"},
			wantErr: "email is required",
		},
		{
			name:    "invalid email format",
			req:     CreateUserRequest{FirstName: "John", LastName: "Doe", Email: "notanemail", Password: "12345678"},
			wantErr: "invalid email format",
		},
		{
			name:    "short password",
			req:     CreateUserRequest{FirstName: "John", LastName: "Doe", Email: "a@b.com", Password: "short"},
			wantErr: "password must be at least 8 characters",
		},
		{
			name:    "password exactly 7 chars",
			req:     CreateUserRequest{FirstName: "John", LastName: "Doe", Email: "a@b.com", Password: "1234567"},
			wantErr: "password must be at least 8 characters",
		},
		{
			name: "valid request",
			req:  CreateUserRequest{FirstName: "John", LastName: "Doe", Email: "john@example.com", Password: "securepassword"},
		},
		{
			name: "password exactly 8 chars",
			req:  CreateUserRequest{FirstName: "John", LastName: "Doe", Email: "john@example.com", Password: "12345678"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.wantErr)
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestLoginRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     LoginRequest
		wantErr string
	}{
		{
			name:    "empty email",
			req:     LoginRequest{Email: "", Password: "password"},
			wantErr: "email is required",
		},
		{
			name:    "empty password",
			req:     LoginRequest{Email: "a@b.com", Password: ""},
			wantErr: "password is required",
		},
		{
			name:    "both empty",
			req:     LoginRequest{Email: "", Password: ""},
			wantErr: "email is required",
		},
		{
			name: "valid request",
			req:  LoginRequest{Email: "a@b.com", Password: "password"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.wantErr)
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestUpdateUserRequest_Validate(t *testing.T) {
	strPtr := func(s string) *string { return &s }

	tests := []struct {
		name    string
		req     UpdateUserRequest
		wantErr string
	}{
		{
			name: "all nil fields is valid",
			req:  UpdateUserRequest{},
		},
		{
			name: "valid email update",
			req:  UpdateUserRequest{Email: strPtr("new@example.com")},
		},
		{
			name:    "invalid email format",
			req:     UpdateUserRequest{Email: strPtr("bademail")},
			wantErr: "invalid email format",
		},
		{
			name:    "short password",
			req:     UpdateUserRequest{Password: strPtr("short")},
			wantErr: "password must be at least 8 characters",
		},
		{
			name: "valid password update",
			req:  UpdateUserRequest{Password: strPtr("longenoughpassword")},
		},
		{
			name: "valid first name update",
			req:  UpdateUserRequest{FirstName: strPtr("Jane")},
		},
		{
			name: "valid last name update",
			req:  UpdateUserRequest{LastName: strPtr("Smith")},
		},
		{
			name: "multiple valid fields",
			req: UpdateUserRequest{
				FirstName: strPtr("Jane"),
				LastName:  strPtr("Smith"),
				Email:     strPtr("jane@smith.com"),
				Password:  strPtr("newpassword123"),
			},
		},
		{
			name:    "valid name but invalid email",
			req:     UpdateUserRequest{FirstName: strPtr("Jane"), Email: strPtr("invalid")},
			wantErr: "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.wantErr)
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

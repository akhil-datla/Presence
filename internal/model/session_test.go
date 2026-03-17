package model

import (
	"testing"
)

func TestCreateSessionRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateSessionRequest
		wantErr string
	}{
		{
			name:    "empty name",
			req:     CreateSessionRequest{Name: ""},
			wantErr: "name is required",
		},
		{
			name: "valid name",
			req:  CreateSessionRequest{Name: "Morning Standup"},
		},
		{
			name: "single character name",
			req:  CreateSessionRequest{Name: "A"},
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

func TestUpdateSessionRequest_Validate(t *testing.T) {
	strPtr := func(s string) *string { return &s }

	tests := []struct {
		name    string
		req     UpdateSessionRequest
		wantErr string
	}{
		{
			name: "nil name is valid",
			req:  UpdateSessionRequest{Name: nil},
		},
		{
			name:    "empty string name",
			req:     UpdateSessionRequest{Name: strPtr("")},
			wantErr: "name cannot be empty",
		},
		{
			name: "valid name",
			req:  UpdateSessionRequest{Name: strPtr("Updated Session")},
		},
		{
			name: "single character name",
			req:  UpdateSessionRequest{Name: strPtr("X")},
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

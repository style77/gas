package accounts

import (
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email       interface{}
		expectedErr bool
	}{
		{"test@example.com", false},
		{"invalid-email", true},
		{"", true},
	}

	for _, tt := range tests {
		err := isValidEmail(tt.email)
		if tt.expectedErr && err == nil {
			t.Errorf("Expected error for email: %v, got nil", tt.email)
		}
		if !tt.expectedErr && err != nil {
			t.Errorf("Unexpected error for email: %v, got %v", tt.email, err)
		}
	}
}

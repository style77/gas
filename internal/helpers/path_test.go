package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPath(t *testing.T) {
	// Retrieve the current user's home directory for comparison
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	// Define test cases
	testCases := []struct {
		name           string
		input          string
		expected       string
		expectingError bool
	}{
		{
			name:           "Path starts with ~ and a slash",
			input:          "~/documents",
			expected:       filepath.Join(homeDir, "documents"),
			expectingError: false,
		},
		{
			name:           "Path is only ~",
			input:          "~",
			expected:       homeDir,
			expectingError: false,
		},
		{
			name:           "Path does not start with ~",
			input:          "/usr/local/bin",
			expected:       "/usr/local/bin",
			expectingError: false,
		},
		{
			name:           "Path starts with ~ but no slash",
			input:          "~documents",
			expected:       filepath.Join(homeDir, "documents"),
			expectingError: false,
		},
		{
			name:           "Empty path",
			input:          "",
			expected:       "",
			expectingError: false,
		},
		{
			name:           "Path starts with multiple tildes",
			input:          "~~/folder",
			expected:       "~~/folder", // Should not be expanded
			expectingError: false,
		},
		{
			name:           "Path with ~ in the middle",
			input:          "/home/~user/documents",
			expected:       "/home/~user/documents",
			expectingError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ExpandPath(tc.input)
			if tc.expectingError {
				if err == nil {
					t.Errorf("Expected error for input '%s', but got none", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("Did not expect an error for input '%s', but got: %v", tc.input, err)
				return
			}
			if result != tc.expected {
				t.Errorf("For input '%s', expected '%s' but got '%s'", tc.input, tc.expected, result)
			}
		})
	}
}

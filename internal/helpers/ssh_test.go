// internal/helpers/helpers_test.go
package helpers

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// MockCommandExecutor is a mock implementation of CommandExecutor.
type MockCommandExecutor struct {
	LookPathFunc    func(file string) (string, error)
	ExecCommandFunc func(name string, arg ...string) *exec.Cmd
}

func (m MockCommandExecutor) LookPath(file string) (string, error) {
	return m.LookPathFunc(file)
}

func (m MockCommandExecutor) ExecCommand(name string, arg ...string) *exec.Cmd {
	return m.ExecCommandFunc(name, arg...)
}

// Helper function to create a temporary file with given content
func createTempFile(t *testing.T, dir, pattern, content string, perm os.FileMode) string {
	t.Helper()
	tmpFile, err := os.CreateTemp(dir, pattern)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	if content != "" {
		if _, err := tmpFile.WriteString(content); err != nil {
			tmpFile.Close()
			t.Fatalf("Failed to write to temp file: %v", err)
		}
	}
	if err := tmpFile.Chmod(perm); err != nil {
		tmpFile.Close()
		t.Fatalf("Failed to set permissions on temp file: %v", err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

func TestIsValidSSHKey(t *testing.T) {
	// Setup: Create temporary files for testing
	validKeyContent := `-----BEGIN RSA PRIVATE KEY-----
MIIBOgIBAAJBAKj34GkxFhD90vcNLYLInFEX6Ppy1tPf9Cnzj4p4WGeKLs1Pt8Qu
KUpRKfFLfRYC9AIKjbJTWit+CqvjWYzvQwECAwEAAQJAIJLixBy2qpFoS4DSmoEm
o3qGy0t6z09AIJtH+5OeRV1be+N4cDYJKffGzDa88vQENZiRm0GRq6a+HPGQMd2k
TQIhAKMSvzIBnni7ot/OSie2TmJLY4SwTQAevXysE2RbFDYdAiEBCUEaRQnMnbp7
9mxDXDf6AU0cN/RPBjb9qSHDcWZHGzUCIG2Es59z8ugGrDY+pxLQnwfotadxd+Uy
v/Ow5T0q5gIJAiEAyS4RaI9YG8EWx/2w0T67ZUVAw8eOMB6BIUg0Xcu+3okCIBOs
/5OiPgoTdSy7bcF9IGpSE8ZgGKzgYQVZeN97YE00
-----END RSA PRIVATE KEY-----`

	invalidKeyContent := "this is not a valid ssh key"

	// Create a valid SSH key file
	validKeyPath := createTempFile(t, "", "valid_key_*", validKeyContent, 0600)
	defer os.Remove(validKeyPath)

	// Create an invalid SSH key file
	invalidKeyPath := createTempFile(t, "", "invalid_key_*", invalidKeyContent, 0600)
	defer os.Remove(invalidKeyPath)

	// Create a file without read permissions
	unreadableKeyPath := createTempFile(t, "", "unreadable_key_*", validKeyContent, 0000)
	defer os.Remove(unreadableKeyPath)

	// Define test cases
	testCases := []struct {
		name          string
		input         interface{}
		setup         func() string // Function to set up any specific conditions
		expectedError error
	}{
		{
			name:          "Non-string input",
			input:         12345,
			expectedError: errors.New("path is not a string"),
		},
		{
			name:          "Valid SSH key path",
			input:         validKeyPath,
			expectedError: nil,
		},
		{
			name:          "Path does not exist",
			input:         "/non/existent/path/id_rsa",
			expectedError: errors.New("invalid path"),
		},
		{
			name:          "Unreadable file",
			input:         unreadableKeyPath,
			expectedError: errors.New("could not read key"),
		},
		{
			name:          "Invalid SSH key content",
			input:         invalidKeyPath,
			expectedError: errors.New("invalid key"),
		},
		{
			name:  "Path starts with ~ and valid key",
			input: "~/.ssh/id_rsa_test_valid",
			setup: func() string {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					t.Fatalf("Failed to get user home directory: %v", err)
				}
				path := filepath.Join(homeDir, ".ssh", "id_rsa_test_valid")
				if err := os.WriteFile(path, []byte(validKeyContent), 0600); err != nil {
					t.Fatalf("Failed to write valid key to %s: %v", path, err)
				}
				return path
			},
			expectedError: nil,
		},
		{
			name:          "Path starts with ~ and does not exist",
			input:         "~/.ssh/non_existent_key",
			expectedError: errors.New("invalid path"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var input interface{}
			if tc.setup != nil {
				path := tc.setup()
				defer os.Remove(path) // Clean up if setup created a file
				input = tc.input
			} else {
				input = tc.input
			}

			err := IsValidSSHKey(input)
			if tc.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error '%v' but got none", tc.expectedError)
				} else if err.Error() != tc.expectedError.Error() {
					t.Errorf("Expected error '%v' but got '%v'", tc.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got '%v'", err)
				}
			}
		})
	}
}

func TestGenerateSSHKey_WithMockExecutor(t *testing.T) {
	mockExecutor := MockCommandExecutor{
		LookPathFunc: func(file string) (string, error) {
			if file == "ssh-keygen" {
				return "/usr/bin/ssh-keygen", nil
			}
			return "", exec.ErrNotFound
		},
		ExecCommandFunc: func(name string, arg ...string) *exec.Cmd {
			// Simulate successful command execution
			// Using "echo" to mimic ssh-keygen behavior without actual key generation
			cmd := exec.Command("echo", "SSH key generated")
			return cmd
		},
	}

	// Setup: Create a temporary directory to act as HOME
	tempHomeDir := t.TempDir()

	// Override the HOME environment variable
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempHomeDir)

	email := "test@example.com"
	keyPath, err := GenerateSSHKey(email, mockExecutor)
	if err != nil {
		t.Errorf("Did not expect an error, but got: %v", err)
	}

	expectedPath := filepath.Join(tempHomeDir, ".ssh", "id_rsa")
	if keyPath != expectedPath {
		t.Errorf("Expected key path '%s', but got '%s'", expectedPath, keyPath)
	}

	// Since we're using "echo", no actual files are created. We can optionally check if the function behaves correctly.
}

func TestGenerateSSHKey_SSHKeygenNotFound(t *testing.T) {
	mockExecutor := MockCommandExecutor{
		LookPathFunc: func(file string) (string, error) {
			return "", exec.ErrNotFound
		},
		ExecCommandFunc: func(name string, arg ...string) *exec.Cmd {
			return exec.Command("echo")
		},
	}

	email := "test@example.com"
	_, err := GenerateSSHKey(email, mockExecutor)
	if err == nil {
		t.Errorf("Expected an error when ssh-keygen is not found, but got none")
	} else {
		expectedErrMsg := "ssh-keygen not found"
		if err.Error() != expectedErrMsg {
			t.Errorf("Expected error '%s', but got '%s'", expectedErrMsg, err.Error())
		}
	}
}

func TestGenerateSSHKey_CommandExecutionFailure(t *testing.T) {
	mockExecutor := MockCommandExecutor{
		LookPathFunc: func(file string) (string, error) {
			if file == "ssh-keygen" {
				return "/usr/bin/ssh-keygen", nil
			}
			return "", exec.ErrNotFound
		},
		ExecCommandFunc: func(name string, arg ...string) *exec.Cmd {
			// Simulate command failure by running a command that exits with status 1
			return exec.Command("false")
		},
	}

	email := "test@example.com"
	_, err := GenerateSSHKey(email, mockExecutor)
	if err == nil {
		t.Errorf("Expected an error due to command failure, but got none")
	} else {
		expectedErrPrefix := "failed to generate ssh key:"
		if !strings.HasPrefix(err.Error(), expectedErrPrefix) {
			t.Errorf("Expected error prefix '%s', but got '%s'", expectedErrPrefix, err.Error())
		}
	}
}

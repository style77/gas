package helpers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"golang.org/x/crypto/ssh"
)

// IsValidSSHKey checks if an ssh key is valid.
func IsValidSSHKey(path interface{}) error {
	filePath, ok := path.(string)
	if !ok {
		return errors.New("path is not a string")
	}

	// Expand the path if it contains '~'
	expandedPath, err := ExpandPath(filePath)
	if err != nil {
		return err
	}

	// Check if the expanded path is valid
	if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
		return errors.New("invalid path")
	}

	// Read the key
	key, err := os.ReadFile(expandedPath)
	if err != nil {
		return errors.New("could not read key")
	}

	// check if key is valid
	_, err = ssh.ParsePrivateKey(key)
	if err != nil {
		return errors.New("invalid key")
	}

	return nil
}

// GenerateSSHKey generates an ssh key.
func GenerateSSHKey(email string) (string, error) {
	var keyPath string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	if runtime.GOOS == "windows" {
		keyPath = filepath.Join(homeDir, ".ssh", "id_rsa")
	} else {
		keyPath = filepath.Join(homeDir, ".ssh", "id_rsa")
	}

	// Ensure the .ssh directory exists
	sshDir := filepath.Dir(keyPath)
	err = os.MkdirAll(sshDir, 0700)
	if err != nil {
		return "", fmt.Errorf("failed to create .ssh directory: %w", err)
	}

	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-C", email, "-f", keyPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Check if ssh-keygen exists
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		return "", errors.New("ssh-keygen not found")
	}

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to generate ssh key: %w", err)
	}

	fmt.Println("SSH key successfully generated at:", keyPath)
	return keyPath, nil
}

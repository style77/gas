package helpers

import (
	"os/user"
	"path/filepath"
	"strings"
)

// ExpandPath expands a path that starts with '~' to the user's home directory.
func ExpandPath(path string) (string, error) {
	// Check if the path starts with '~' and expand it
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		// Replace '~' with the user's home directory
		return filepath.Join(usr.HomeDir, path[1:]), nil
	}
	return path, nil
}

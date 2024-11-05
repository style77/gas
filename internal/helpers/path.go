// helpers.go
package helpers

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expands a path that starts with '~' to the user's home directory.
func ExpandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") && !strings.HasPrefix(path, "~~") {
		usr, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(usr, path[1:]), nil
	}
	return path, nil
}

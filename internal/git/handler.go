package git

import (
	"os"
	"os/exec"
	"strings"
)

// HandleGitCommand runs the provided git command with the provided arguments.
func HandleGitCommand(args []string) error {
	gitCmd := exec.Command("git", args...)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	gitCmd.Stdin = os.Stdin

	err := gitCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// UpdateGlobalGitConfig updates the global git configuration with the provided username and email.
func UpdateGlobalGitConfig(username, email string) {
	exec.Command("git", "config", "--global", "user.name", username).Run()
	exec.Command("git", "config", "--global", "user.email", email).Run()
}

func IsCurrentGlobal(email string) bool {
	currentEmail, _ := exec.Command("git", "config", "--global", "user.email").Output()

	return strings.TrimSpace(string(currentEmail)) == string(email)
}

func GetCurrentGlobal() string {
	currentEmail, _ := exec.Command("git", "config", "--global", "user.email").Output()

	return strings.TrimSpace(string(currentEmail))
}

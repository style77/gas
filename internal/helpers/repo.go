package helpers

import (
	"fmt"
	"regexp"
)

// ExtractUserAndRepo extracts the user and repo name from a remote URL
func ExtractUserAndRepo(remoteUrl string) (string, string, error) {
	// Match patterns like "git@hostname:user/repo.git"
	sshPattern := regexp.MustCompile(`git@[\w.-]+:([\w.-]+)/([\w.-]+)\.git`)

	// Match patterns like "https://hostname/user/repo.git"
	httpsPattern := regexp.MustCompile(`https://[\w.-]+/([\w.-]+)/([\w.-]+)\.git`)

	var matches []string
	if sshPattern.MatchString(remoteUrl) {
		matches = sshPattern.FindStringSubmatch(remoteUrl)
	} else if httpsPattern.MatchString(remoteUrl) {
		matches = httpsPattern.FindStringSubmatch(remoteUrl)
	} else {
		return "", "", fmt.Errorf("unsupported remote URL format: %s", remoteUrl)
	}

	if len(matches) == 3 {
		return matches[1], matches[2], nil
	}
	return "", "", fmt.Errorf("failed to parse remote URL: %s", remoteUrl)
}

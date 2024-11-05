package git

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type GitHubClient interface {
	FetchPublicKeys(username string) ([]string, error)
	IsGithubUsernameValid(username string) error
}

type RealGitHubClient struct{}

// fetchGitHubPublicKeys fetches the public keys of a github user.
func (c *RealGitHubClient) FetchPublicKeys(username string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/keys", username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("could not fetch public keys")
	}

	var keys []struct {
		Key string `json:"key"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return nil, err
	}

	publicKeys := make([]string, len(keys))
	for i, key := range keys {
		publicKeys[i] = key.Key
	}

	return publicKeys, nil
}

// isGithubUsernameValid checks if a username exists on github.
func (c *RealGitHubClient) IsGithubUsernameValid(username string) error {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	resp, err := http.Get(url)
	if err != nil {
		return errors.New("could not check if username is valid")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return errors.New("username not found")
	}

	return nil
}

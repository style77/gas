package repo

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/style77/gas/internal/accounts"
	"github.com/style77/gas/internal/git"
	"github.com/style77/gas/internal/helpers"
)

func SetRemoteUrl(account *accounts.Account, remoteName string) error {
	remoteUrl := git.GetCurrentRemoteUrl(remoteName) // todo consider other remotes

	user, repo, err := helpers.ExtractUserAndRepo(remoteUrl)
	if err != nil {
		return err
	}

	newRemoteUrl := fmt.Sprintf("git@%s:%s/%s.git", account.SSHAlias, user, repo)

	var isCorrectRemoteURL bool
	survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Do you want to set the remote URL to '%s'?", newRemoteUrl),
	}, &isCorrectRemoteURL)

	if !isCorrectRemoteURL {
		return errors.New("remote URL not set")
	}

	err = git.SetRemoteUrl("origin", newRemoteUrl)
	if err != nil {
		return err
	}

	return nil
}

package accounts

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/style77/gas/internal/git"
	"github.com/style77/gas/internal/helpers"
	"golang.org/x/crypto/ssh"
)

var interactiveAddAccountInvestigationQuestions = []*survey.Question{
	{
		Name:     "Email",
		Prompt:   &survey.Input{Message: "What is email address associated with the github account?"},
		Validate: isValidEmail,
	},
	{
		Name:     "Name",
		Prompt:   &survey.Input{Message: "What is the github name you wish to use? It might be your github account username (\"johnDoe98\") or your real name (\"John Doe\")."},
		Validate: survey.Required,
	},
}

// InteractiveNewAccount prompts the user for information to add a new account.
func InteractiveNewAccount() {
	investigationAnswers := struct {
		Email string
		Name  string
	}{}

	err := survey.Ask(interactiveAddAccountInvestigationQuestions, &investigationAnswers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	githubClient := &git.RealGitHubClient{}

	err, isExistingGithubAccount := validateAndPromptForName(&investigationAnswers.Name, githubClient.IsGithubUsernameValid)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !isExistingGithubAccount {
		fmt.Println("You have chosen to use a name that is not a valid GitHub username. Keep in mind GAS won't be able to verify SSH keys for this account.")
	}

	var SSHKeyExists bool
	err = survey.AskOne(&survey.Confirm{Message: "Do you have an ssh key you would like to use with this account? If you don't have one, you will be prompted through the process of creating one."}, &SSHKeyExists)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var SSHKeyPath string
	if SSHKeyExists {
		err = survey.AskOne(&survey.Input{Message: "What is the path to the ssh key you would like to use with this account?"}, &SSHKeyPath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = helpers.IsValidSSHKey(SSHKeyPath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if isExistingGithubAccount {
			isValid, err := isValidSSHKeyForGitHub(SSHKeyPath, investigationAnswers.Name, githubClient)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if !isValid {
				fmt.Println("The key you provided is not associated with the account you are trying to add.")
				return
			}
		} else {
			fmt.Println("Since the name you provided is not a valid GitHub username, GAS cannot verify the key you provided. Continuing with the account creation process.")
		}
	} else {
		SSHKeyPath, err = helpers.GenerateSSHKey(investigationAnswers.Email)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	sshAlias := handleSSHConfig(SSHKeyPath)

	account := Account{
		Email:      investigationAnswers.Email,
		Name:       investigationAnswers.Name,
		SSHKeyPath: SSHKeyPath,
		SSHAlias:   sshAlias,
	}
	SaveAccountToConfig(account)
}

// handleSSHConfig handles the SSH configuration for the provided SSH key path.
func handleSSHConfig(sshKeyPath string) string {
	var sshAlias string
	sshConfigPath := filepath.Join(os.Getenv("HOME"), ".ssh", "config")

	configFile, err := os.OpenFile(sshConfigPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		fmt.Printf("Could not open SSH config file: %v\n", err)
		return ""
	}

	var configContent strings.Builder
	_, err = io.Copy(&configContent, configFile)
	if err != nil {
		fmt.Printf("Could not read SSH config file: %v\n", err)
		configFile.Close()
		return ""
	}
	configFile.Close()

	existingAlias := findExistingAlias(configContent.String(), sshKeyPath)
	if existingAlias != "" {
		sshAlias = existingAlias
		fmt.Printf("Using existing SSH alias: %s\n", sshAlias)
	} else {
		err = survey.AskOne(&survey.Input{Message: "Enter a unique alias for this SSH key (e.g., github-work):"}, &sshAlias, survey.WithValidator(survey.Required))
		if err != nil {
			fmt.Println(err.Error())
			return ""
		}

		newConfigEntry := fmt.Sprintf(`
Host %s
    HostName github.com
    User git
    IdentityFile %s
`, sshAlias, sshKeyPath)

		if runtime.GOOS == "windows" {
			newConfigEntry = strings.ReplaceAll(newConfigEntry, "\\", "/")
		}

		configFile, err = os.OpenFile(sshConfigPath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Printf("Could not open SSH config file for appending: %v\n", err)
			return ""
		}
		defer configFile.Close()

		_, err = configFile.WriteString(newConfigEntry)
		if err != nil {
			fmt.Printf("Could not write to SSH config file: %v\n", err)
			return ""
		}

		fmt.Printf("Added SSH configuration for alias '%s'.\n", sshAlias)
	}

	return sshAlias
}

// findExistingAlias is a Helper function to find if an alias already exists for the given ssh key path in the SSH config file.
func findExistingAlias(configContent, sshKeyPath string) string {
	lines := strings.Split(configContent, "\n")
	currentAlias := ""

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Could not determine user home directory: %v\n", err)
		return ""
	}

	normalizedSSHKeyPath := sshKeyPath
	if strings.HasPrefix(sshKeyPath, "~/") {
		normalizedSSHKeyPath = filepath.Join(homeDir, sshKeyPath[2:])
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Host ") {
			currentAlias = strings.Fields(trimmed)[1]
		}
		if strings.HasPrefix(trimmed, "IdentityFile") {
			identityFilePath := strings.Fields(trimmed)[1]

			if strings.HasPrefix(identityFilePath, "~") {
				identityFilePath = strings.Replace(identityFilePath, "~", homeDir, 1)
			}

			if identityFilePath == normalizedSSHKeyPath {
				return currentAlias
			}
		}
	}

	return ""
}

// isValidUsername checks if a username is valid based on the rules of github usernames.
func isValidEmail(email interface{}) error {
	if survey.Required(email) != nil {
		return fmt.Errorf("email is required")
	}

	emailStr := email.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(emailStr) {
		return fmt.Errorf("invalid email")
	}

	return nil
}

// validateAndPromptForName checks if the name is valid and prompts the user if necessary.
func validateAndPromptForName(name *string, validateFunc func(string) error) (error, bool) {
	if err := isNameValid(*name, validateFunc); err != nil {
		var wantsToUseName bool
		err := survey.AskOne(&survey.Confirm{Message: *name + " looks like a Github username, but GAS found out that it might be non-existent or invalid. Would you like to use this name anyway?", Default: true}, &wantsToUseName)
		if err != nil {
			return err, false
		}

		if !wantsToUseName {
			// Prompt for a new name recursively
			return promptForNewName(name, validateFunc), false
		}
	}
	return nil, true
}

// promptForNewName prompts the user for a new name and checks validity recursively.
func promptForNewName(name *string, validateFunc func(string) error) error {
	err := survey.AskOne(&survey.Input{Message: "What is the GitHub name you wish to use? It might be your GitHub account username (\"johnDoe98\") or your real name (\"John Doe\")."}, name)
	if err != nil {
		return err
	}

	// Check the new name validity
	err, _ = validateAndPromptForName(name, validateFunc)
	return err
}

// isNameValid checks if a username exists on github.
func isNameValid(name string, validateFunc func(string) error) error {
	if regexp.MustCompile(`^[a-zA-Z0-9._%+-]+$`).MatchString(name) {
		return validateFunc(name)
	}

	return nil
}

// isValidSSHKeyForGitHub checks if an ssh key is valid for a github account.
func isValidSSHKeyForGitHub(filePath string, username string, client git.GitHubClient) (bool, error) {
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	privateKey, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return false, err
	}

	publicKey := privateKey.PublicKey()
	publicKeyString := string(ssh.MarshalAuthorizedKey(publicKey))

	publicKeys, err := client.FetchPublicKeys(username)
	if err != nil {
		return false, err
	}

	for _, key := range publicKeys {
		if strings.TrimSpace(key) == strings.TrimSpace(publicKeyString) {
			return true, nil
		}
	}

	return false, nil
}

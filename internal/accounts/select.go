package accounts

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/style77/gas/internal/git"
)

func InteractiveSelectAccount() (Account, error) {
	accounts := GetAccounts()

	if len(accounts) == 0 {
		return Account{}, fmt.Errorf("no accounts found")
	}

	accountNames := []string{}
	for _, account := range accounts {
		accountName := account.Name
		if git.IsCurrentGlobal(account.Email) {
			accountName += " (global)"
		}

		accountNames = append(accountNames, accountName)
	}

	prompt := &survey.Select{
		Message: "Select the account you would like to use:",
		Options: accountNames,
	}

	var selectedAccountName string
	survey.AskOne(prompt, &selectedAccountName)

	selectedAccountName = strings.Replace(selectedAccountName, " (global)", "", -1)
	selectedAccount, _ := GetAccount(selectedAccountName)

	return selectedAccount, nil
}

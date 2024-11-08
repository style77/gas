package accounts

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/viper"
	"github.com/style77/gas/internal/git"
)

type Account struct {
	Name       string
	Email      string
	SSHKeyPath string
	SSHAlias   string
	Id         int
}

// SaveAccountToConfig saves the account information to the configuration file.
func SaveAccountToConfig(account Account) {
	accounts := viper.GetStringMap("accounts")

	// check if account already exists
	if _, ok := accounts[account.Name]; ok {
		overwrite := false
		err := survey.AskOne(&survey.Confirm{
			Message: "Do you want to overwrite the existing account?",
		}, &overwrite)
		if err != nil {
			fmt.Println(err)
			return
		}

		if !overwrite {
			fmt.Printf("Account '%s' already exists. Exiting.\n", account.Name)
			return
		}
	}

	maxID := 0
	for _, acc := range accounts {
		accMap, ok := acc.(map[string]interface{})
		if ok {
			if idValue, ok := accMap["id"]; ok {
				if id, valid := idValue.(int); valid && id > maxID {
					maxID = id
				}
			}
		}
	}
	newAccountID := maxID + 1

	account.Id = newAccountID

	accounts[account.Name] = map[string]interface{}{
		"name":       account.Name,
		"email":      account.Email,
		"sshkeypath": account.SSHKeyPath,
		"sshalias":   account.SSHAlias,
		"id":         account.Id,
	}

	viper.Set("accounts", accounts)

	err := viper.WriteConfig()
	if err != nil {
		fmt.Printf("Failed to save account '%s'. Error: %s\n", account.Name, err)
		return
	}

	fmt.Printf("Account '%s' added successfully.\n", account.Name)
}

func GetAccount(name string) (Account, error) {
	accounts := viper.GetStringMap("accounts")
	account, ok := accounts[name]
	if !ok {
		return Account{}, fmt.Errorf("account '%s' not found", name)
	}

	accountMap, ok := account.(map[string]interface{})
	if !ok {
		return Account{}, fmt.Errorf("account '%s' has an invalid format", name)
	}

	return Account{
		Name:       accountMap["name"].(string),
		Email:      accountMap["email"].(string),
		SSHKeyPath: accountMap["sshkeypath"].(string),
		SSHAlias:   accountMap["sshalias"].(string),
		Id:         accountMap["id"].(int),
	}, nil
}

func GetAccounts() []Account {
	accounts := viper.GetStringMap("accounts")
	var result []Account

	for key, account := range accounts {
		accountMap, ok := account.(map[string]interface{})
		if !ok {
			fmt.Printf("Account '%s' has an invalid format or type is incorrect\n", key)
			continue
		}

		name, nameOk := accountMap["name"].(string)
		email, emailOk := accountMap["email"].(string)
		sshKeyPath, sshKeyOk := accountMap["sshkeypath"].(string)

		if !nameOk || !emailOk || !sshKeyOk {
			fmt.Printf("Account '%s' is missing required fields or has invalid types\n", key)
			continue
		}

		result = append(result, Account{
			Name:       name,
			Email:      email,
			SSHKeyPath: sshKeyPath,
		})
	}

	return result
}

func (a *Account) String() string {
	return fmt.Sprintf("Name: %s, Email: %s", a.Name, a.Email)
}

func (a *Account) Delete() {
	accounts := viper.GetStringMap("accounts")
	delete(accounts, a.Name)
	viper.Set("accounts", accounts)

	err := viper.WriteConfig()
	if err != nil {
		fmt.Printf("Failed to delete account '%s'. Error: %s\n", a.Name, err)
		return
	}

	fmt.Printf("Account '%s' deleted successfully.\n", a.Name)
}

func (a *Account) SetGlobal() {
	git.UpdateGlobalGitConfig(a.Name, a.Email)
}

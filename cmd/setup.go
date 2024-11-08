package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/style77/gas/internal/accounts"
	"github.com/style77/gas/internal/repo"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		accountRaw, _ := cmd.Flags().GetString("account")

		var account accounts.Account
		if accountRaw == "" {
			var err error
			account, err = accounts.InteractiveSelectAccount()
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			var err error
			account, err = accounts.GetAccount(accountRaw)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		if account.Name == "" {
			fmt.Println("No account selected.")
			return
		}

		remoteName, _ := cmd.Flags().GetString("remoteName")
		err := repo.SetRemoteUrl(&account, remoteName)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Configured repo's remote URL to use account '%s'.", account.Name)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	setupCmd.Flags().StringP("account", "a", "", "Account to switch to. This should be the name of the account you wish to switch to.")
	setupCmd.Flags().StringP("remoteName", "r", "origin", "Remote name to set the URL for.")
}

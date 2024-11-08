package cmd

import (
	"github.com/spf13/cobra"
	"github.com/style77/gas/internal/accounts"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Add a new GitHub account",
	Long:  `Add a new GitHub account to the list of accounts on this machine. Run this command to interactively provide the account details.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if no arguments are passed, run the interactive version of the command
		if len(args) == 0 {
			accounts.InteractiveNewAccount()
		} else {
			// TODO implement non-interactive version of adding
			accounts.InteractiveNewAccount()
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addAccountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addAccountCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

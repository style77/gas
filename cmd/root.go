package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/style77/gas/internal/git"
)

var rootCmd = &cobra.Command{
	Use:   "gas",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// isUnknownCommandError checks if the error is an "unknown command" error.
func isUnknownCommandError(err error) bool {
	return err != nil && len(err.Error()) > 0 && err.Error()[0:15] == "unknown command"
}

func Execute() {
	rootCmd.SilenceErrors = true

	rootErr := rootCmd.Execute()
	if rootErr != nil {
		if isUnknownCommandError(rootErr) {
			args := os.Args[1:]

			// Check if this is desired account
			currentAccount := git.GetCurrentGlobal()

			var isProperAccount bool
			survey.AskOne(&survey.Confirm{
				Message: fmt.Sprintf("Do you want to run '%s' as '%s'?", strings.Join(args, " "), currentAccount),
			}, &isProperAccount)

			if isProperAccount {
				err := git.HandleGitCommand(args)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Exiting.")
				os.Exit(1)
			}
		} else {
			rootCmd.SilenceErrors = false
			fmt.Println(rootErr)
			os.Exit(1)
		}
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".gas")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create it
			err := viper.WriteConfigAs(home + "/.gas.yaml")
			if err != nil {
				fmt.Println("Failed to create config file: ", err)
			}
		} else {
			// Config file was found but another error was produced
			fmt.Println("Failed to read config file: ", err)
		}
	}
}

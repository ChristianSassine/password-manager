package cmd

import (
	"os"

	"github.com/ChristianSassine/password-manager/pass-cli/internal/output"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pass-cli",
	Short: "CLI to manage passwords",
	Long: `pass-cli is a CLI that communicates with the password manager server.
It allows to quickly manage and read passwords from the CLI. 
The user and password should be set as environment variables under 'PASS_USERNAME' and 'PASS_PASSWORD' respectively.
e.g: PASS_USERNAME=Hello PASS_PASSWORD=World`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func getBool(cmd *cobra.Command, name string) bool {
	result, err := cmd.Flags().GetBool(name)
	if err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
	return result
}

func getInt(cmd *cobra.Command, name string) int {
	result, err := cmd.Flags().GetInt(name)
	if err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
	return result
}

func getString(cmd *cobra.Command, name string) string {
	result, err := cmd.Flags().GetString(name)
	if err != nil {
		output.Error("%v", err)
		os.Exit(1)
	}
	return result
}

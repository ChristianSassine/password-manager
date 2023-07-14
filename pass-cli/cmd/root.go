package cmd

import (
	"os"

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

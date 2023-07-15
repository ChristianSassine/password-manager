package cmd

import (
	"os"

	"github.com/ChristianSassine/password-manager/pass-cli/internal/manager"
	"github.com/ChristianSassine/password-manager/pass-cli/internal/output"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [key]",
	Short: "Adds a new password with the key provided",
	Long: `This command requests the server to generate a new key for the current user.
It will take the key as argument.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		returnPassword, err := cmd.Flags().GetBool("output")
		if err != nil {
			output.Error("%v", err)
			os.Exit(1)
		}

		isQuiet, err := cmd.Flags().GetBool("quiet")
		if err != nil {
			output.Error("%v", err)
			os.Exit(1)
		}

		output.SetOutput(isQuiet)
		manager.AddPassword(args[0], returnPassword)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("output", "o", false, "Will display the password")
	addCmd.Flags().BoolP("quiet", "q", false, "Will silence descriptions and only return the result (if enabled)")
}

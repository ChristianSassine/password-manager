package cmd

import (
	"os"

	"github.com/ChristianSassine/password-manager/pass-cli/internal/manager"
	"github.com/ChristianSassine/password-manager/pass-cli/internal/output"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [key]",
	Short: "Removes the password that matches the key.",
	Long: `This command sends a request to the server to remove the password matching with the key.
It will take the password's key as argument.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		isQuiet, err := cmd.Flags().GetBool("quiet")
		if err != nil {
			output.Error("%v", err)
			os.Exit(1)
		}

		output.SetOutput(isQuiet)
		manager.RemovePassword(args[0])
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().BoolP("quiet", "q", false, "Will silence descriptions")
}

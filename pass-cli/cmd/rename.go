package cmd

import (
	"os"

	"github.com/ChristianSassine/password-manager/pass-cli/internal/manager"
	"github.com/ChristianSassine/password-manager/pass-cli/internal/output"
	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename [key] [newKeyName]",
	Short: "Renames a password's key.",
	Long: `This command sends a request to the server to rename a password's key with the [newKeyName].
It will take the password's current [key] and the [newKeyName] as arguments.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			cmd.Help()
		}

		isQuiet, err := cmd.Flags().GetBool("quiet")
		if err != nil {
			output.Error("%v", err)
			os.Exit(1)
		}

		output.SetOutput(isQuiet)
		key, oldKey := args[0], args[1]
		manager.RenameKey(key, oldKey)
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
	renameCmd.Flags().BoolP("quiet", "q", false, "Will silence descriptions")
}

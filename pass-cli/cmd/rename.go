package cmd

import (
	"os"

	"github.com/ChristianSassine/password-manager/pass-cli/internal/manager"
	"github.com/ChristianSassine/password-manager/pass-cli/internal/output"
	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename [key] [newKeyName]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

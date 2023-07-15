package cmd

import (
	"os"

	"github.com/ChristianSassine/password-manager/pass-cli/internal/manager"
	"github.com/ChristianSassine/password-manager/pass-cli/internal/output"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [key]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

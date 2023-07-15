package cmd

import (
	"github.com/ChristianSassine/password-manager/pass-cli/internal/manager"
	"github.com/ChristianSassine/password-manager/pass-cli/internal/output"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieves the password that matchs the key",
	Long: `This command sends a request to the server to retrieve the password matching with the key.
It will take the key as argument.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		returnPassword, err := cmd.Flags().GetBool("output")
		output.Print("returned password %v", returnPassword)
		if err != nil {
			output.Error("%v", err)
		}
		isQuiet, err := cmd.Flags().GetBool("quiet")
		output.Print("isQuiet %v", isQuiet)
		if err != nil {
			output.Error("%v", err)
		}
		output.SetOutput(isQuiet)
		manager.GetPassword(args[0], returnPassword)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().BoolP("output", "o", false, "Will display the password")
	getCmd.Flags().BoolP("quiet", "q", false, "Will silence descriptions and only return the result (if enabled)")
}

package cmd

import (
	"github.com/ChristianSassine/password-manager/pass-cli/internal/manager"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Sets the server's URL and creates a user",
	Long: `This command is used to set the server's URL which will be saved in a config file.
	This command is also used to create an account.
	Refer to the flags for this command.`,
	Run: func(cmd *cobra.Command, args []string) {

		url, _ := cmd.Flags().GetString("url")
		creds, _ := cmd.Flags().GetString("creds")
		if url != "" {
			manager.SetURL(url)
		}

		if creds != "" {
			manager.CreateAccount(creds)
		}

		if creds == "" && url == "" {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringP("url", "u", "", "Sets the servers URL (exclude the schema). e.g: -u=127.0.0.1:4200")
	configCmd.Flags().StringP("creds", "c", "", "Requests the server to create an account <username>:<password>. e.g: -a=hello:world")
}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ChristianSassine/password-manager/pass-cli/internal/manager"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieves the password that matchs the key",
	Long: `This command sends a request to the server to retrieve the password matching with the key.
	It will take the key as argument.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Handle flags
		if len(args) == 0 {
			cmd.Help()
		}
		manager.GetPassword(args[0])
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().BoolP("output", "o", false, "Will display the password")
	getCmd.Flags().BoolP("quiet", "q", false, "Will silence descriptions and only return the result (if enabled)")
}

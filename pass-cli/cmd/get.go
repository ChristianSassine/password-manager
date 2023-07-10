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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

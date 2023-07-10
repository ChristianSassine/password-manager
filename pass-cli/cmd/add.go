/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ChristianSassine/password-manager/pass-cli/internal/manager"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [key]",
	Short: "Adds a new password with the key provided",
	Long: `This command requests the server to generate a new key for the current user.
It will take the key as argument.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Handle flags
		if len(args) == 0 {
			cmd.Help()
		}
		manager.AddPassword(args[0])
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("output", "o", false, "Will display the password")
	addCmd.Flags().BoolP("quiet", "q", false, "Will silence descriptions and only return the result (if enabled)")
}

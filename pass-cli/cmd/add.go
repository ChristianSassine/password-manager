package cmd

import (
	"github.com/ChristianSassine/password-manager/pass-cli/internal/manager"
	"github.com/ChristianSassine/password-manager/pass-cli/internal/output"
	"github.com/ChristianSassine/password-manager/pass-cli/internal/requests"
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

		returnPassword := getBool(cmd, "output")
		isQuiet := getBool(cmd, "quiet")

		// Parameters
		length := getInt(cmd, "characters")
		disableLowerLetters := getBool(cmd, "lower")
		disableUpperLetters := getBool(cmd, "upper")
		disableDigits := getBool(cmd, "digits")
		disableSpecial := getBool(cmd, "special")

		output.SetOutput(isQuiet)
		manager.AddPassword(requests.Parameters{
			Key: args[0], Length: length, LowerLetters: !disableLowerLetters,
			UpperLetters: !disableUpperLetters, Digits: !disableDigits,
			Symbols: !disableSpecial,
		},
			returnPassword)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP("output", "o", false, "Will display the password")
	addCmd.Flags().BoolP("quiet", "q", false, "Will silence descriptions and only return the result (if enabled)")

	// Password config flags
	addCmd.Flags().IntP("characters", "c", 32, "choose how many characters should be included in the password (by default 32). eg: -c 10")
	addCmd.Flags().BoolP("lower", "l", false, "generate a password without lowerCase letters")
	addCmd.Flags().BoolP("upper", "u", false, "generate a password without UpperCase letters")
	addCmd.Flags().BoolP("digits", "d", false, "generate a password without digits")
	addCmd.Flags().BoolP("special", "s", false, "generate a password without special characters")
}

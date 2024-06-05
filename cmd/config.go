package cmd

import "github.com/spf13/cobra"

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Validate the configuration file",
		RunE: func(_ *cobra.Command, _ []string) error {
			// TODO: Implement the config command

			// get the --init flag value
			// init, _ := cmd.Flags().GetBool("init")
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(configCmd)

	// Add --init flag to the config command to create a new default config file & --validate flag to validate the existing config file
	configCmd.Flags().BoolP("init", "i", false, "Create a new default configuration file")
	configCmd.Flags().BoolP("validate", "v", false, "Validate the existing configuration file")
}

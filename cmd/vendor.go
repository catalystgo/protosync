package cmd

import "github.com/spf13/cobra"

var (
	vendorCmd = &cobra.Command{
		Use:   "vendor",
		Short: "Download proto files from a remote repository",
		Long:  `Download proto files from a remote repository`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	rootCmd.AddCommand(vendorCmd)
}

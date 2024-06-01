package cmd

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use:   "protosync",
		Short: "protosync is a tool to sync proto files from a remote repository",
		Long:  `protosync is a tool to sync proto files from a remote repository`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

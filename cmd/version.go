package cmd

import (
	"fmt"

	"github.com/catalystgo/protosync/internal/build"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:     "version",
		Short:   "Print the version number of protosync",
		Long:    `Print the version number of protosync`,
		Aliases: []string{"ver"},
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(build.Version)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

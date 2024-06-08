package cmd

import (
	"github.com/catalystgo/protosync/internal/build"
	"github.com/catalystgo/xro-log/log"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:     "version",
		Short:   "Print the version number of protosync",
		Long:    `Print the version number of protosync`,
		Aliases: []string{"ver"},
		Run: func(cmd *cobra.Command, _ []string) {
			formatType, err := cmd.Flags().GetString("format")
			if err != nil {
				log.Fatalf("read output flag: %v", err)
			}

			log.Debugf("format type is %s", formatType)

			err = svc.PrintVersion(build.Version, build.Commit, build.Date, formatType)
			if err != nil {
				log.Fatalf("print version: %v", err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().String("format", "text", "output format")
}

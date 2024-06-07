package cmd

import (
	"github.com/catalystgo/protosync/internal/build"
	"github.com/catalystgo/xro-log/log"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of protosync",
		Long:  `Print the version number of protosync`,
		Run: func(cmd *cobra.Command, _ []string) {
			outputType, err := cmd.Flags().GetString("output")
			if err != nil {
				log.Fatalf("read output flag: %v", err)
			}

			log.Debugf("output type is %s", outputType)

			err = svc.PrintVersion(build.Version, build.Commit, build.Date, outputType)
			if err != nil {
				log.Fatalf("print version: %v", err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().StringP("output", "o", "", "output format (json)")
}

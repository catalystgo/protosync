package cmd

import (
	"github.com/catalystgo/protosync/internal/config"
	"github.com/catalystgo/xro-log/log"
	"github.com/spf13/cobra"
)

var (
	vendorCmd = &cobra.Command{
		Use:     "vendor",
		Short:   "Download proto files from a remote repository",
		Long:    `Download proto files from a remote repository`,
		Aliases: []string{"ven"},
		Run: func(_ *cobra.Command, _ []string) {
			c := config.Get()
			for _, d := range c.Dependencies {
				if err := svc.Download(d.Source, c.OutDir); err != nil {
					log.Warnf("download %s => %v", d.Source, err)
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(vendorCmd)
}

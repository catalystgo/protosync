package cmd

import (
	"github.com/catalystgo/protosync/internal/config"
	"github.com/spf13/cobra"
)

var (
	vendorCmd = &cobra.Command{
		Use:   "vendor",
		Short: "Download proto files from a remote repository",
		Long:  `Download proto files from a remote repository`,
		RunE: func(_ *cobra.Command, _ []string) error {
			c := config.Get()

			for _, d := range c.Dependencies {
				if err := svc.Download(d.Source, c.OutDir); err != nil {
					return err
				}
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(vendorCmd)
}

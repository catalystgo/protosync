package cmd

import (
	log "github.com/catalystgo/logger/cli"
	"github.com/spf13/cobra"
)

var (
	outputDir string
)

var (
	vendorCmd = &cobra.Command{
		Use:     "vendor",
		Short:   "Download proto files from a remote repository",
		Long:    `Download proto files from a remote repository`,
		Aliases: []string{"ven"},
		Run: func(_ *cobra.Command, _ []string) {
			for _, d := range cfg.Dependencies {
				for _, f := range d.Sources {
					if err := svc.Download(f, cfg.AbsOutDir, d.Path); err != nil {
						log.Errorf("download %s => %v", f, err)
					}
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(vendorCmd)

	vendorCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "", "output directory path")
}

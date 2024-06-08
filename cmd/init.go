package cmd

import (
	"github.com/catalystgo/xro-log/log"
	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:     "init",
		Short:   "Generate default configuration file",
		Long:    "Generate default configuration file",
		Aliases: []string{"i"},
		Run: func(_ *cobra.Command, _ []string) {
			if err := svc.GenConfig(configPath); err != nil {
				log.Fatalf("generate default configuration file: %v", err)
			}
			log.Infof("configuration file created")
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}

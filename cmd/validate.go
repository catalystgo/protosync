package cmd

import (
	"github.com/catalystgo/protosync/internal/config"
	"github.com/catalystgo/xro-log/log"
	"github.com/spf13/cobra"
)

var (
	validateCmd = &cobra.Command{
		Use:     "validate",
		Short:   "Validate the configuration file",
		Long:    "Validate the configuration file",
		Aliases: []string{"val"},
		Run: func(cmd *cobra.Command, _ []string) {
			if _, err := config.Load(configPath, outputDir); err != nil {
				log.Fatalf("validate configuration file: %v", err)
			}
			log.Info("configuration file is valid")
		},
	}
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

package cmd

import (
	"github.com/catalystgo/protosync/internal/config"
	"github.com/catalystgo/xro-log/log"
	"github.com/spf13/cobra"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Init or Check the configuration file",
		Run: func(cmd *cobra.Command, _ []string) {

			// get the --init flag value
			if init, _ := cmd.Flags().GetBool("init"); init {
				if err := svc.GenConfig(configPath); err != nil {
					log.Fatalf("generate configuration file: %v", err)
				}
				log.Infof("configuration file created")
			}

			// get the --check flag value
			if check, _ := cmd.Flags().GetBool("check"); check {
				if _, err := config.Load(configPath); err != nil {
					log.Fatalf("load configuration file: %v", err)
				}
				log.Info("configuration file is valid")
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(configCmd)

	// --init flag to the config command to create a new default config file
	configCmd.Flags().BoolP("init", "i", false, "Create a new default configuration file")

	// --Check flag to Check the existing config file
	configCmd.Flags().BoolP("check", "c", false, "Check the existing configuration file")
}

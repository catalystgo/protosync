package cmd

import (
	"github.com/catalystgo/logger/log"
	"github.com/catalystgo/protosync/internal/config"
	"github.com/catalystgo/protosync/internal/domain"
	"github.com/catalystgo/protosync/internal/downloader"
	"github.com/catalystgo/protosync/internal/http"
	"github.com/catalystgo/protosync/internal/service"
	"github.com/spf13/cobra"
)

var (
	cfg *config.Config

	configPath string
	verbose    bool
)

var (
	httpClient = http.NewClient()

	// Downloaders

	githubClient    = downloader.NewGithub(httpClient)
	gitlabClient    = downloader.NewGitlab(httpClient)
	bitbucketClient = downloader.NewBitbucket(httpClient)

	// Writer

	writer = service.NewWriteProvider()

	// Services

	svc = service.New(writer)
)

var (
	rootCmd = &cobra.Command{
		Use:          "protosync",
		Short:        "protosync is a tool to sync proto files from a remote repository",
		Long:         `protosync is a tool to sync proto files from a remote repository`,
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			if verbose {
				log.SetLevel(log.LevelDebug)
			}

			// Skip loading config if the command is not vendor
			// since vendor command is the only command that requires config
			if cmd.Name() != vendorCmd.Name() {
				return
			}

			log.Debugf("pre run on command: %s", cmd.CommandPath())
			log.Debugf("loading config file: %s", configPath)

			// Load config file
			c, err := config.Load(configPath, outputDir)
			if err != nil {
				log.Fatalf("load config: %v\n", err)
			}

			log.Debug("config loaded")
			log.Debugf("config output directory: %s", c.Directory)

			cfg = c // set config to global variable

			// Register downloaders for each external domain
			for _, d := range c.Domains {
				downloader, ok := svc.GetDownloader(d.API)

				log.Debugf("registering downloader for domain %s with api %s", d.Host, d.API)

				if !ok {
					log.Fatalf("missing downloader for domain %s", d.Host)
				}

				svc.Register(d.Host, downloader)
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "file", "f", "protosync.yml", "config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	svc.Register(domain.DefaultDomainGithub, githubClient)
	svc.Register(domain.DefaultDomainGitlab, gitlabClient)
	svc.Register(domain.DefaultDomainBitbucket, bitbucketClient)
}

func Execute() error {
	return rootCmd.Execute()
}

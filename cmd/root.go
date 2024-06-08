package cmd

import (
	"github.com/catalystgo/protosync/internal/config"
	"github.com/catalystgo/protosync/internal/domain"
	"github.com/catalystgo/protosync/internal/downloader"
	"github.com/catalystgo/protosync/internal/http"
	"github.com/catalystgo/protosync/internal/service"
	"github.com/catalystgo/xro-log/log"
	"github.com/spf13/cobra"
)

var (
	configPath string
	verbose    bool
	outputDir  string
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
			log.Debug("loading config")

			// Load config file
			c, err := config.Load(configPath, outputDir)
			if err != nil {
				log.Fatalf("load config: %v", err)
			}

			log.Debug("config loaded")
			log.Debugf("config output directory: %s", c.OutDir)

			// Register downloaders for each external domain
			for _, d := range c.Domains {
				downloader, ok := svc.GetDownloader(d.API)

				log.Debugf("registering downloader for domain %s with api %s", d.Name, d.API)

				if !ok {
					log.Fatalf("missing downloader for domain %s", d.Name)
				}

				svc.Register(d.Name, downloader)
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "file", "f", "protosync.yml", "config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "", "output directory path")

	svc.Register(domain.DefaultDomainGithub, githubClient)
	svc.Register(domain.DefaultDomainGitlab, gitlabClient)
	svc.Register(domain.DefaultDomainBitbucket, bitbucketClient)
}

func Execute() error {
	return rootCmd.Execute()
}

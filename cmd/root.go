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
			c, err := config.Load(configPath)
			if err != nil {
				log.Fatalf("load config: %v", err)
			}

			log.Debug("config loaded")

			// Register downloaders for each external domain
			for _, d := range c.Domains {
				apiDomain := domain.GetAPIDomain(d.API)

				log.Debugf("registering downloader for domain %s with api %s", d.Domain, apiDomain)

				switch apiDomain {
				case domain.DefaultDomainGithub:
					svc.Register(d.Domain, githubClient)
				case domain.DefaultDomainGitlab:
					svc.Register(d.Domain, gitlabClient)
				case domain.DefaultDomainBitbucket:
					svc.Register(d.Domain, bitbucketClient)
				default:
					log.Fatalf("missing downloader for domain %s", d.Domain)
				}
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config-file", "f", "protosync.yml", "config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	svc.Register(domain.DefaultDomainGithub, githubClient)
	svc.Register(domain.DefaultDomainGitlab, gitlabClient)
	svc.Register(domain.DefaultDomainBitbucket, bitbucketClient)
}

func Execute() error {
	return rootCmd.Execute()
}

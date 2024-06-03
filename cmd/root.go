package cmd

import (
	"fmt"

	"github.com/catalystgo/protosync/internal/config"
	"github.com/catalystgo/protosync/internal/domain"
	"github.com/catalystgo/protosync/internal/downloader"
	"github.com/catalystgo/protosync/internal/http"
	"github.com/catalystgo/protosync/internal/service"
	"github.com/spf13/cobra"
)

var configPath string

var (
	httpClient = http.NewClient()

	// Downloaders

	githubClient    = downloader.NewGithub(httpClient)
	gitlabClient    = downloader.NewGitlab(httpClient)
	bitbucketClient = downloader.NewBitbucket(httpClient)

	// Services

	svc = service.New()
)

var (
	rootCmd = &cobra.Command{
		Use:          "protosync",
		Short:        "protosync is a tool to sync proto files from a remote repository",
		Long:         `protosync is a tool to sync proto files from a remote repository`,
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Load config file
			c, err := config.Load(configPath)
			if err != nil {
				return err
			}

			// Register downloaders for each external domain
			for _, d := range c.Domains {
				apiDomain := domain.GetAPIDomain(d.Api)

				switch apiDomain {
				case domain.DefaultDomainGithub:
					svc.Register(d.Domain, githubClient)
				case domain.DefaultDomainGitlab:
					svc.Register(d.Domain, gitlabClient)
				case domain.DefaultDomainBitbucket:
					svc.Register(d.Domain, bitbucketClient)
				default:
					return fmt.Errorf("missing downloader for domain %s", d.Domain)
				}
			}

			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config-file", "f", "proto-sync.yml", "config file")

	svc.Register(domain.DefaultDomainGithub, githubClient)
	svc.Register(domain.DefaultDomainGitlab, gitlabClient)
	svc.Register(domain.DefaultDomainBitbucket, bitbucketClient)
}

func Execute() error {
	return rootCmd.Execute()
}

package downloader

const (
	bitbucketDefaultDomain = "bitbucket.org"
	githubDefaultDomain    = "github.com"
	gitlabDefaultDomain    = "gitlab.com"
)

type Downloader interface {
	GetFile(domain string, path string) ([]byte, error)
}

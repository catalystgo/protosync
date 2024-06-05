package downloader

import (
	"fmt"

	"github.com/catalystgo/protosync/internal/domain"
)

type Github struct {
	client httpClient
}

func NewGithub(httpClient httpClient) *Github {
	return &Github{
		client: httpClient,
	}
}

func (g *Github) GetFile(f *domain.File) ([]byte, error) {
	return getFile(g.client, g.getUrl(f), f)
}

// getUrl returns the URL for the file
func (g *Github) getUrl(f *domain.File) string {
	return fmt.Sprintf("https://%s/%s/%s/blob/%s/%s?raw=true",
		f.Domain, f.User, f.Repo, f.Ref, f.Path,
	)
}

package downloader

import (
	"fmt"

	"github.com/catalystgo/protosync/internal/domain"
)

type Gitlab struct {
	client httpClient
}

func NewGitlab(httpClient httpClient) *Gitlab {
	return &Gitlab{
		client: httpClient,
	}
}

// GetFile downloads the file from the Gitlab repository
func (g *Gitlab) GetFile(f *domain.File) ([]byte, error) {
	return getFile(g.client, g.getURL(f), f)
}

// getURL returns the URL for the file
func (g *Gitlab) getURL(f *domain.File) string {
	return fmt.Sprintf("https://%s/%s/%s/-/raw/%s/%s",
		f.Domain, f.User, f.Repo, f.Ref, f.Path,
	)
}

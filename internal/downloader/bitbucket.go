package downloader

import (
	"fmt"

	"github.com/catalystgo/protosync/internal/domain"
)

type Bitbucket struct {
	client httpClient
}

func NewBitbucket(httpClient httpClient) *Bitbucket {
	return &Bitbucket{
		client: httpClient,
	}
}

func (b *Bitbucket) GetFile(f *domain.File) ([]byte, error) {
	return getFile(b.client, b.getURL(f), f)
}

// getURL returns the URL for the file
func (b *Bitbucket) getURL(f *domain.File) string {
	return fmt.Sprintf("https://%s/%s/%s/raw/%s/%s",
		f.Domain, f.User, f.Repo, f.Ref, f.Path,
	)
}

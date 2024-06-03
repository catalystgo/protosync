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
	return nil, fmt.Errorf("bitbucket client not implemented")
}

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

func (g *Gitlab) GetFile(f *domain.File) ([]byte, error) {
	return nil, fmt.Errorf("gitlab client not implemented")
}

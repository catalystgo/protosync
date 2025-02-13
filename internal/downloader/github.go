package downloader

import (
	"fmt"

	"github.com/catalystgo/protosync/internal/domain"
)

type Github struct {
	client client
}

func NewGithub(client client) *Github {
	return &Github{
		client: client,
	}
}

func (g *Github) GetFile(f *domain.File) ([]byte, error) {
	return getFile(g.client, g.getURL(f), f)
}

func (g *Github) getURL(f *domain.File) string {
	return fmt.Sprintf("https://%s/%s/%s/blob/%s/%s?raw=true",
		f.Domain, f.User, f.Repo, f.Ref, f.Path,
	)
}

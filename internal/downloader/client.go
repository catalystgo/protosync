package downloader

import (
	"fmt"
	"net/http"
	"path"

	log "github.com/catalystgo/logger/cli"
	"github.com/catalystgo/protosync/internal/domain"
)

type client interface {
	Get(url string, opts ...func(*http.Request)) ([]byte, error)
}

func getFile(cli client, url string, f *domain.File) ([]byte, error) {
	log.Debugf("fetch content from: %s", url)

	content, err := cli.Get(url)
	if err != nil {
		return nil, err
	}

	filePath := fmt.Sprintf("%s@%s", path.Join(f.Domain, f.User, f.Repo, f.Path), f.Ref)
	log.Infof("got content for: %s", filePath)

	return content, nil
}

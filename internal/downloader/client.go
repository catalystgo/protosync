package downloader

import (
	"fmt"
	"net/http"
	"path"

	"github.com/catalystgo/protosync/internal/domain"
	"github.com/catalystgo/xro-log/log"
)

// func WithAuthToken(token string) func(*http.Request) {
// 	return func(req *http.Request) {
// 		req.Header.Set("Authorization", "Bearer "+token)
// 	}
// }

type httpClient interface {
	Get(url string, opts ...func(*http.Request)) ([]byte, error)
}

func getFile(client httpClient, url string, f *domain.File) ([]byte, error) {
	log.Debugf("downloading content from: %s", url)

	content, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	filePath := fmt.Sprintf("%s@%s", path.Join(f.Domain, f.User, f.Repo, f.Path), f.Ref)
	log.Infof("downloaded content for: %s", filePath)

	return content, nil
}

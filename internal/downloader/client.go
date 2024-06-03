package downloader

import "net/http"

func WithAuthToken(token string) func(*http.Request) {
	return func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}

type httpClient interface {
	Get(url string, opts ...func(*http.Request)) ([]byte, error)
}

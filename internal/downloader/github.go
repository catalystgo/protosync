package downloader

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/catalystgo/protosync/internal/domain"
)

var (
	ErrContentNotFound = fmt.Errorf("content not found in response")
	ErrContentInvalid  = func(content interface{}) error {
		return fmt.Errorf("invalid content type: %T", content)
	}
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
	u := g.getUrl(f)

	b, err := g.client.Get(u)
	if err != nil {
		return nil, err
	}

	content, err := getContent(b)
	if err != nil {
		return nil, err
	}

	// Decode base64 content
	b, err = decodeBase64(content)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (g *Github) getUrl(f *domain.File) string {
	return fmt.Sprintf("https://api.%s/repos/%s/%s/contents/%s?ref=%s",
		f.Domain, f.User, f.Repo, f.Path, f.Ref,
	)
}

func decodeBase64(s string) ([]byte, error) {
	var out strings.Builder

	lines := strings.Split(s, "\n")
	for _, line := range lines {
		b, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			// continue
			return nil, fmt.Errorf("decode base64: %s => %w", line, err)
		}

		// write to out
		if _, err := out.Write(b); err != nil {
			return nil, err
		}
	}

	return []byte(out.String()), nil
}

// getContent extracts the content from the response body.

func getContent(b []byte) (string, error) {
	// Using Unmarshal to parse the JSON response on a struct
	// didn't work for some reason, so I'm using a map instead.
	// TODO: Use a struct to parse the JSON response
	var jsonResp map[string]interface{}
	if err := json.Unmarshal(b, &jsonResp); err != nil {
		return "", err
	}

	contentObj, ok := jsonResp["content"]
	if !ok || contentObj == nil {
		return "", ErrContentNotFound
	}

	content, ok := contentObj.(string)
	if !ok {
		return "", ErrContentInvalid(contentObj)
	}

	return content, nil
}

package domain

import (
	"fmt"
	"strings"

	log "github.com/catalystgo/logger/cli"
)

var (
	ErrInvalidFile = func(file string, message string) error {
		return fmt.Errorf("invalid source format: %s => %s", file, message)
	}
)

type File struct {
	Domain string
	User   string
	Repo   string
	Path   string
	Ref    string
}

func ParseFile(file string) (*File, error) {
	log.Debugf("parsing file link: %s", file)

	// Get the file's domain
	slashIdx := strings.Index(file, "/")
	if slashIdx == -1 {
		return nil, ErrInvalidFile(file, "missing domain")
	}

	domain := file[:slashIdx]

	if domain == "" {
		return nil, ErrInvalidFile(file, "domain is empty")
	}

	// Get the file's ref
	atIdx := strings.Index(file, "@")
	if atIdx == -1 {
		return nil, ErrInvalidFile(file, "missing ref")
	}

	ref := file[atIdx+1:]

	if ref == "" {
		return nil, ErrInvalidFile(file, "ref is empty")
	}

	// Get the file's path
	// Example:
	// - github.com/catalystgo/protosync/internal/domain/file.go => internal/domain/file.go
	// - gitlab.com/catalystgo/protosync/internal/domain/file.go => internal/domain/file.go
	// - bitbucket.org/catalystgo/protosync/internal/domain/file.go => internal/domain/file.go
	path := file[slashIdx+1 : atIdx]

	if path == "" {
		return nil, ErrInvalidFile(file, "path is empty")
	}

	pathParts := strings.Split(path, "/")
	if len(pathParts) < 3 {
		return nil, ErrInvalidFile(file, "path is invalid")
	}

	user := pathParts[0]
	repo := pathParts[1]
	path = strings.Join(pathParts[2:], "/") // Full file path under the repo

	if user == "" {
		return nil, ErrInvalidFile(file, "user is empty")
	}
	if repo == "" {
		return nil, ErrInvalidFile(file, "repo is empty")
	}
	if !strings.HasSuffix(path, ".proto") {
		return nil, ErrInvalidFile(file, "only .proto extension is allowed")
	}

	f := &File{
		Domain: domain,
		User:   user,
		Repo:   repo,
		Path:   path,
		Ref:    ref,
	}

	return f, nil
}

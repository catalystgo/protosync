package service

import (
	_ "embed"
	"fmt"
	"path"
	"sync"

	"github.com/catalystgo/helpers"
	"github.com/catalystgo/protosync/internal/domain"
)

var (
	ErrInvalidFile = func(file string, message string) error {
		return fmt.Errorf("invalid file format: %s => %s", file, message)
	}
)

type Downloader interface {
	GetFile(file *domain.File) ([]byte, error)
}

type Service struct {
	mu sync.RWMutex

	downloaders map[string]Downloader
}

func New() *Service {
	return &Service{downloaders: make(map[string]Downloader)}
}

func (s *Service) Register(domain string, d Downloader) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.downloaders[domain] = d
}

func (s *Service) GetDownloader(domain string) (Downloader, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.downloaders[domain]
	return d, ok
}

func (s *Service) Download(file string, outDirectory string, outPath string) error {
	f, err := domain.ParseFile(file)
	if err != nil {
		return err
	}

	// Get downloader by domain

	s.mu.RLock()
	d, ok := s.downloaders[f.Domain]
	s.mu.RUnlock()

	if !ok {
		return ErrInvalidFile(file, "unknown domain")
	}

	// Download file content

	content, err := d.GetFile(f)
	if err != nil {
		return err
	}

	// Set output path. If not provided, use the default path
	// based on the file domain, user, repo, and path values
	// from the file.
	if outPath == "" {
		outPath = path.Join(outDirectory, f.Domain, f.User, f.Repo, f.Path)
	} else {
		outPath = path.Join(outDirectory, outPath, path.Base(f.Path))
	}

	// Write file content

	err = helpers.SaveFile(outPath, content, &helpers.SaveFileOpt{Override: true})
	if err != nil {
		return err
	}

	return nil
}

var (
	//go:embed template/protosync.yml
	configContent string

	configName = "protosync.yml"
	configPath = path.Join(".", configName)
)

// GenConfig generates a default configuration file
// for the protosync tool in the current directory
func (s *Service) GenConfig(configFile string) error {
	if configFile == "" {
		configFile = configPath
	}
	return helpers.SaveFile(configFile, []byte(configContent), &helpers.SaveFileOpt{Override: false})
}

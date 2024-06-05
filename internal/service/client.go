package service

import (
	_ "embed"
	"fmt"
	"path"
	"sync"

	"github.com/catalystgo/protosync/internal/domain"
)

var (
	ErrorInvalidFile = func(file string, message string) error {
		return fmt.Errorf("invalid file format: %s => %s", file, message)
	}
)

type Downloader interface {
	GetFile(file *domain.File) ([]byte, error)
}

type Writer interface {
	Write(file string, content []byte) error
}

type Service struct {
	mu sync.RWMutex

	writer      Writer
	downloaders map[string]Downloader
}

func New(writer Writer) *Service {
	return &Service{
		writer:      writer,
		downloaders: make(map[string]Downloader),
	}
}

func (s *Service) Register(domain string, d Downloader) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.downloaders[domain] = d
}

func (s *Service) Download(file string, outDir string) error {
	f, err := domain.ParseFile(file)
	if err != nil {
		return err
	}

	// Get downloader by domain

	s.mu.RLock()
	d, ok := s.downloaders[f.Domain]
	s.mu.RUnlock()

	if !ok {
		return ErrorInvalidFile(file, "unknown domain")
	}

	// Download file content

	content, err := d.GetFile(f)
	if err != nil {
		return err
	}

	// Write file content

	outPath := path.Join(outDir, f.Domain, f.User, f.Repo, f.Path)
	err = s.writer.Write(outPath, content)
	if err != nil {
		return err
	}

	return nil
}

var (
	//go:embed template/proto-sync.yml
	configContent string

	configName = "proto-sync.yml"
	configPath = path.Join(".", configName)
)

// GenConfig generates a default configuration file
// for the protosync tool in the current directory
func (s *Service) GenConfig() error {
	return s.writer.Write(configPath, []byte(configContent))
}

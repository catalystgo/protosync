package service

import (
	"fmt"
	"os"
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

type Service struct {
	mu          sync.RWMutex
	downloaders map[string]Downloader
}

func New() *Service {
	return &Service{
		downloaders: make(map[string]Downloader),
	}
}

func (s *Service) Register(name string, d Downloader) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.downloaders[name] = d
}

func (s *Service) DownloadAll(outDir string, files ...string) error {
	for _, file := range files {
		err := s.Download(file, outDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Download(file string, outDir string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	f, err := domain.ParseFile(file)
	if err != nil {
		return err
	}

	// Get downloader by domain

	d, ok := s.downloaders[f.Domain]
	if !ok {
		return ErrorInvalidFile(file, "unknown domain")
	}

	// Download file content

	content, err := d.GetFile(f)
	if err != nil {
		return err
	}

	// Create output file and write content

	outPath := path.Join(outDir, f.Domain, f.Path)

	// Create output directory if not exists
	err = os.MkdirAll(path.Dir(outPath), os.ModePerm)
	if err != nil {
		return err
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}

	defer func() { _ = outFile.Close() }()

	_, err = outFile.Write(content)
	if err != nil {
		return err
	}

	return nil
}

package service

import (
	"os"
	"path"

	"github.com/catalystgo/xro-log/log"
)

type WriteProvider struct{}

func NewWriteProvider() *WriteProvider {
	return &WriteProvider{}
}

func (p *WriteProvider) Write(file string, content []byte) error {
	// Create the directory if it doesn't exist
	err := os.MkdirAll(path.Dir(file), os.ModePerm)
	if err != nil {
		return err
	}

	if _, err := os.Stat(file); err == nil {
		log.Warnf("overwriting existing file: %s", file)
	}

	// Create the file
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer func() { _ = f.Close() }()

	// Write the content to the file
	_, err = f.Write(content)
	if err != nil {
		return err
	}

	return nil
}

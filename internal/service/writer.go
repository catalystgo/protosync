package service

import (
	"os"
	"path"
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

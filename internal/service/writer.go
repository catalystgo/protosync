package service

import (
	"os"
	"path"

	log "github.com/catalystgo/logger/cli"
)

type writer struct{}

func newWriter() *writer {
	return &writer{}
}

func (p *writer) Write(file string, content []byte, overide bool) error {
	// Create the directory if it doesn't exist
	err := os.MkdirAll(path.Dir(file), os.ModePerm)
	if err != nil {
		return err
	}

	if _, err := os.Stat(file); err == nil {
		if !overide {
			log.Warnf("config file [%s] already exists so skipping", file)
			return nil
		}
		log.Debugf("overide existing file: %s", file)
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

	log.Debugf("content written to file: %s", file)

	return nil
}

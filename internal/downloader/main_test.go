package downloader

import (
	"testing"

	log "github.com/catalystgo/logger/cli"
)

func TestMain(m *testing.M) {
	log.SetLevel(log.LevelWarn)
	m.Run()
}

package downloader

import (
	"testing"

	"github.com/catalystgo/logger/log"
)

func TestMain(m *testing.M) {
	log.SetLevel(log.LevelWarn)
	m.Run()
}

package downloader

import (
	"testing"

	"github.com/catalystgo/xro-log/log"
)

func TestMain(m *testing.M) {
	log.SetLevel(log.LevelWarn)
	m.Run()
}

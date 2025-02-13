package build

import (
	"runtime/debug"
	"strings"
)

var (
	Version = "DEV"
)

func init() {
	if strings.ToLower(Version) == "dev" {
		info, ok := debug.ReadBuildInfo()
		if ok && info.Main.Version != "" {
			Version = info.Main.Version
		}
	}
}

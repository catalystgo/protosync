package build

import (
	"runtime/debug"
)

var (
	Version = "DEV"
)

func init() {
	if Version == "DEV" {
		info, ok := debug.ReadBuildInfo()
		if ok && info.Main.Version != "" {
			Version = info.Main.Version
		}
	}
}

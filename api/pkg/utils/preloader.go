package utils

import (
	"strings"
	"time"
)

func Preload() (time.Time, string) {
	start := time.Now()
	line := strings.Repeat("-", 15)
	PrintStd(Std, "", "\n%s [service] initializing %s\n", line, line)
	return start, line
}

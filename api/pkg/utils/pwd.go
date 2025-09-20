package utils

import (
	"log"
	"os"
	"path/filepath"
)

const rootDetector = "go.mod"

// Root finds and returns the root directory of the project by locating the go.mod file.
func Root() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Walk up the directory tree until we find go.mod
	for {
		if _, err = os.Stat(filepath.Join(dir, rootDetector)); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			log.Fatal("failed to detect the project root path")
		}

		dir = parent
	}
}

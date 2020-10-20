package utility

import (
	"os"
)

// MakeDirectory : Create the directory, if absent
func MakeDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModeDir)
	}
}

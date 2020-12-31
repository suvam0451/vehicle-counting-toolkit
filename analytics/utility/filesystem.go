package utility

import (
	"fmt"
	"os"
)

// MakeDirectory : Create the directory, if absent
func MakeDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, os.ModeDir)
	}
}

// EnsureFile makes sure the file targeted exists
func EnsureFile(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer file.Close()
	}
}

// PathExists returns whether the given file or directory exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return false, err
}
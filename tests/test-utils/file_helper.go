package testutils

import (
	"errors"
	"os"
)

func RemoveFileIfExists(filePath string) error {
	if _, err := os.Stat(filePath); !errors.Is(err, os.ErrNotExist) {
		err = os.Remove(filePath)
		return err
	}
	return nil
}

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

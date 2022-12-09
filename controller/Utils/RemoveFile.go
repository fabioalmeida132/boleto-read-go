package Utils

import (
	"github.com/pkg/errors"
	"os"
)

func RemoveFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return errors.New("Error removing file")
	}
	return nil
}

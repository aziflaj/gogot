package fileutils

import (
	"os"
	"strings"
)

func CurrentRef() (string, error) {
	headFile, err := os.OpenFile(HeadFilePath, os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}

  ref := strings.Split(FileContents(headFile), ": ")[1]
	
  return ref, nil
}

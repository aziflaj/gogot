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

func CurrentBranch() (string, error) {
	ref, err := CurrentRef()
	if err != nil {
		return "", err
	}

	splitRef := strings.Split(ref, "/")

	return splitRef[len(splitRef)-1], nil
}

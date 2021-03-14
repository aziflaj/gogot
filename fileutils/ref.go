package fileutils

import (
	"io/ioutil"
	"strings"
)

func CurrentRef() (string, error) {
	content, err := ioutil.ReadFile(HeadFilePath)
	if err != nil {
		return "", err
	}

	ref := strings.Split(string(content), ": ")[1]
	return ref, nil
}

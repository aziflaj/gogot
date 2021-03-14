package core

import (
	"io/ioutil"
	"strings"

	"github.com/aziflaj/gogot/files"
)

func CurrentRef() (string, error) {
	content, err := ioutil.ReadFile(files.HeadFilePath)
	if err != nil {
		return "", err
	}

	ref := strings.Split(string(content), ": ")[1]
	return ref, nil
}

package fileutils

import (
	"fmt"
	"io/ioutil"
	"os"
)

func CurrentBranchCommitsFile() (*os.File, error) {
	currentBranchPath, err := CurrentRef()
	if err != nil {
		return nil, err
	}

	commitsFile, err := os.Open(fmt.Sprintf("%s/%s", GogotDir, currentBranchPath))
	if err != nil {
		return nil, err
	}

	return commitsFile, nil
}

func CreateAndOpenCommitFile(hash string) (file *os.File, err error) {
	objectDirPath := fmt.Sprintf("%s/%s", ObjectsDir, hash[0:2])
	os.Mkdir(objectDirPath, 0755)
	objectPath := fmt.Sprintf("%s/%s", objectDirPath, hash[2:])
	file, err = os.OpenFile(objectPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return
}

func ReadCommitContents(hash string) (string, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s", ObjectsDir, hash[0:2], hash[2:]))
	if err != nil {
		return "", nil
	}

	return string(content), nil
}

func ReadBlobContents(hash string) ([]byte, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s", ObjectsDir, hash[0:2], hash[2:]))
	if err != nil {
		return []byte{}, nil
	}

	return content, nil
}

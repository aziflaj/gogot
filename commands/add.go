package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aziflaj/gogot/core"
	"github.com/aziflaj/gogot/fileutils"
)

// Add ...
func Add(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: gogot add <path1> [<path2>] ...")
		os.Exit(1)
	}

	for _, path := range args {
		filesInPath, err := fileutils.AllPaths(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, file := range filesInPath {
			addFile(file)
		}
	}
}

func addFile(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hash := core.HashBytes(content)
	blob := core.CompressBytes(content)

	blobDir := fmt.Sprintf("%s/%s", fileutils.ObjectsDir, hash[0:2])
	os.Mkdir(blobDir, 0755)

	blobPath := fmt.Sprintf("%s/%s", blobDir, hash[2:])
	createBlobFile(blobPath, blob)

	appendToIndexFile(hash, path)
}

func createBlobFile(path string, content []byte) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Some error occurred while creating blob for " + path)
		os.Exit(1)
	}
	defer file.Close()

	file.Write(content)
}

func appendToIndexFile(hash string, path string) {
	f, err := os.OpenFile(fileutils.IndexFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	found := false
	indexedPaths := fileutils.ReadLines(f)

	for idx, hashIndex := range indexedPaths {
		hashAndPath := strings.Split(hashIndex, " ")
		fmt.Println(hashAndPath)
		if hashAndPath[1] == path {
			// already in index, update
			indexedPaths[idx] = fmt.Sprintf("%s %s", hash, path)
			found = true
		}
	}

	if !found {
		indexedPaths = append(indexedPaths, fmt.Sprintf("%s %s", hash, path))
	}

	f.Seek(0, 0)

	for _, hashIndex := range indexedPaths {
		if _, err := f.WriteString(hashIndex + "\n"); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
}

package commands

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aziflaj/gogot/core"
	"github.com/aziflaj/gogot/fileutils"
)

// Add ...
func Add(args []string) {
	if len(args) < 1 {
		log.Fatal("Usage: gogot add <path1> [<path2>] ...")
	}

	for _, path := range args {
		filesInPath, err := fileutils.AllPaths(path)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range filesInPath {
			addFile(file)
		}
	}
}

func addFile(path string) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	content := fileutils.FileBytes(file)

	hash := core.HashBytes(content)
	blob := core.CompressBytes(content)

	blobDir := fmt.Sprintf("%s/%s", fileutils.ObjectsDir, hash[0:2])
	os.Mkdir(blobDir, 0755)
	blobPath := fmt.Sprintf("%s/%s", blobDir, hash[2:])
	createBlobFile(blobPath, blob)

	appendToIndexFile(hash, path)
}

func createBlobFile(path string, content string) {
	file, err := os.Create(path)
	if err != nil {
		log.Printf("Some error occurred while creating blob for %s", path)
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString(content)
}

func appendToIndexFile(hash string, path string) {
	f, err := os.OpenFile(fileutils.IndexFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	found := false
	indexedPaths := fileutils.ReadLines(f)

	for idx, hashIndex := range indexedPaths {
		hashAndPath := strings.Split(hashIndex, " ")
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
			log.Fatal(err)
		}
	}
}

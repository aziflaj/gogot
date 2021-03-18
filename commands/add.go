package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/aziflaj/gogot/core"
	"github.com/aziflaj/gogot/fileutils"
)

var patterns = ignoredPatterns()

// Add ...
func Add(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: gogot add <path1> [<path2>] ...")
		os.Exit(1)
	}

	for _, filepath := range args {
		addRecursive(filepath)
	}
}

func addRecursive(filepath string) {
	for _, pattern := range patterns {
		if match, _ := path.Match(pattern, filepath); match {
			return
		}
	}

	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		fmt.Printf("File doesn't exist: %v\n", filepath)
		fmt.Println(err)
		os.Exit(1)
	}

	if info.IsDir() {
		// files, _ := os.ReadDir(".") // TODO: Once upgrading to Go 1.16
		files, _ := ioutil.ReadDir(filepath)
		for _, file := range files {
			if file.Name() == fileutils.GogotDir {
				continue
			}
			addRecursive(fmt.Sprintf("%s/%s", filepath, file.Name()))
		}
	} else {
		addFile(filepath)
	}
}

func addFile(path string) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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

func ignoredPatterns() []string {
	objectFile, err := os.Open(fileutils.GogotIgnore)
	if err != nil {
		return []string{}
	}

	return fileutils.ReadLines(objectFile)
}

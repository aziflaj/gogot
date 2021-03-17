package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

	appendToIndexFile(fmt.Sprintf("%s %s\n", hash, path))
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

func appendToIndexFile(index string) {
	f, err := os.OpenFile(fileutils.IndexFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	if _, err := f.WriteString(index); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

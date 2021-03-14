package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/aziflaj/gogot/core"
)

// TODO: move to somewhere else
const (
	gogotDir   = ".gogot"
	objectsDir = ".gogot/objects"
	indexPath  = ".gogot/index"
	headPath   = ".gogot/HEAD"

	gogotIgnore = ".gogotignore"
)

// Add ...
// TODO: Can't add a single file. Make a fix
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
	patterns := ignoredPatterns()
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
		files, _ := ioutil.ReadDir(filepath)
		for _, file := range files {
			if file.Name() == gogotDir {
				continue
			}
			addRecursive(fmt.Sprintf("%s/%s", filepath, file.Name()))
		}
	} else {
		addFile(filepath)
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

	blobDir := fmt.Sprintf("%s/%s", objectsDir, hash[0:2])
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
	f, err := os.OpenFile(core.IndexFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	if _, err := f.WriteString(index); err != nil {
		log.Println(err)
	}
}

func ignoredPatterns() (paths []string) {
	objectFile, err := os.Open(core.GogotIgnore)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(objectFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		paths = append(paths, scanner.Text())
	}
	objectFile.Close()

	return
}

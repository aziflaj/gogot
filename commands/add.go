package commands

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
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
func Add(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: gogot add [FILE1] [FILE2] ...")
		os.Exit(1)
	}

	if _, err := os.Stat(gogotDir); os.IsNotExist(err) {
		fmt.Println("Not a Gogot repository")
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

	sha := hashContent(content)
	blob := compressContent(content)

	blobDir := fmt.Sprintf("%s/%s", objectsDir, sha[0:2])
	os.Mkdir(blobDir, 0755)

	blobPath := fmt.Sprintf("%s/%s", blobDir, sha[2:])
	createBlobFile(blobPath, blob)

	appendToIndexFile(fmt.Sprintf("%s %s\n", sha, path))
}

func hashContent(content []byte) string {
	hasher := sha1.New()
	hasher.Write(content)
	sha1Bytes := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return string(sha1Bytes)
}

func compressContent(content []byte) []byte {
	var buffer bytes.Buffer

	writer := zlib.NewWriter(&buffer)
	writer.Write(content)
	writer.Close()

	return buffer.Bytes()
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
	f, err := os.OpenFile(indexPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	if _, err := f.WriteString(index); err != nil {
		log.Println(err)
	}
}

func ignoredPatterns() (paths []string) {
	objectFile, err := os.Open(gogotIgnore)
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

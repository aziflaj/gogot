package fileutils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// Memoizing this because of recursive AllPaths
var GogotIgnorePatterns = IgnoredPatterns()

func FileBytes(file *os.File) (content []byte) {
	scanner := bufio.NewScanner(file)
	return scanner.Bytes()
}

func FileContents(file *os.File) (content string) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		content += "\n" + scanner.Text()
	}
	return
}

func ReadLines(file *os.File) (lines []string) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return
}

func AllPaths(filepath string) (paths []string, err error) {
	for _, pattern := range GogotIgnorePatterns {
		if match, _ := path.Match(pattern, filepath); match {
			return nil, nil
		}
	}

	info, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		files, _ := ioutil.ReadDir(filepath)
		for _, file := range files {
			dirPaths, err := AllPaths(fmt.Sprintf("%s/%s", filepath, file.Name()))
			if err != nil {
				return nil, err
			}
			paths = append(paths, dirPaths...)
		}
	} else {
		paths = append(paths, filepath)
	}

	return paths, nil
}

func IgnoredPatterns() (ignored []string) {
	ignored = make([]string, 2)
	ignored[0] = fmt.Sprintf("./%s", GogotDir)
	ignored[1] = fmt.Sprintf("./%s", GogotIgnore)

	objectFile, err := os.Open(GogotIgnore)
	if err != nil {
		return []string{}
	}

	// I think I like this ... more than ES6's ...
	return append(ignored, ReadLines(objectFile)...)
}

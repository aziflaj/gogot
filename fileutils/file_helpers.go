package fileutils

import (
	"bufio"
	"fmt"
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

func checkPathIgnored(filepath string) (bool, error) {
	for _, pattern := range GogotIgnorePatterns {
		match, err := path.Match(pattern, filepath)
		if err != nil {
			return false, err
		}

		if match {
			return match, nil
		}
	}

	return false, nil
}

func AllPaths(filepath string) (paths []string, err error) {
	isIgnored, err := checkPathIgnored(filepath)
	if err != nil || isIgnored {
		return nil, err
	}

	info, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		files, _ := os.ReadDir(filepath)
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

	objectFile, err := os.Open(ignored[1])
	if err != nil {
		return ignored
	}

	// I think I like this ... more than ES6's ...
	return append(ignored, ReadLines(objectFile)...)
}

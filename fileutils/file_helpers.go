package fileutils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func ReadLines(file *os.File) (lines []string) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return
}

func AllPaths(filepath string) (paths []string, err error) {
	patterns := IgnoredPatterns() // can be constant?
	for _, pattern := range patterns {
		if match, _ := path.Match(pattern, filepath); match {
			return nil, nil
		}
	}

	if filepath == GogotDir {
		return nil, nil
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

func IgnoredPatterns() []string {
	objectFile, err := os.Open(GogotIgnore)
	if err != nil {
		return []string{}
	}

	return ReadLines(objectFile)
}

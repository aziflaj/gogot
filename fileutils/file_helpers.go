package fileutils

import (
	"bufio"
	"os"
)

func ReadLines(file *os.File) (paths []string) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		paths = append(paths, scanner.Text())
	}
	return
}

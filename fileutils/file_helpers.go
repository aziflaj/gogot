package fileutils

import (
	"bufio"
	"os"
)

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

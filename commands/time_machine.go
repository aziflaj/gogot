package commands

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"strings"
)

// TimeMachine ...
func TimeMachine(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: gogot time-machine [COMMIT-ID] [FILE-PATH]")
		os.Exit(1)
	}

	commitID := args[0]
	filePath := args[1]

	treeHash, _, _ := readCommitObjectContent(commitID)
	readObject(treeHash, filePath)
}

func readObject(treeHash string, filePath string) {
	result := readObjectContent(treeHash)

	scanner := bufio.NewScanner(strings.NewReader(result))
	for scanner.Scan() {
		splitContent := strings.Split(scanner.Text(), " ")
		objectType := splitContent[0]
		hash := splitContent[1]
		objectName := splitContent[2]
		if objectType == "blob" && filePath == objectName {
			// found the file
			blob := readBlobContent(hash)
			result = decompressContent(blob)
			fmt.Println(result)
		} else if objectType == "tree" && strings.HasPrefix(filePath, objectName) {
			// found the dir
			pathWithoutPrefix := strings.Split(filePath, objectName+"/")[1]
			readObject(hash, pathWithoutPrefix)
		}
	}
}

func decompressContent(content []byte) string {
	buffer := bytes.NewBuffer(content)
	sb := new(strings.Builder)

	reader, err := zlib.NewReader(buffer)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	io.Copy(sb, reader)

	return sb.String()
}

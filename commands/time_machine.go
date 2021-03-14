package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aziflaj/gogot/core"
)

// TimeMachine ...
func TimeMachine(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: gogot time-machine <commit-id> <file-path>")
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
			result = core.DecompressBytes(blob)
			fmt.Println(result)
		} else if objectType == "tree" && strings.HasPrefix(filePath, objectName) {
			// found the dir
			pathWithoutPrefix := strings.Split(filePath, objectName+"/")[1]
			readObject(hash, pathWithoutPrefix)
		}
	}
}

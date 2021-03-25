package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aziflaj/gogot/fileutils"
)

// Log ...
func Log(args []string) {
	commitsFile, err := fileutils.CurrentBranchCommitsFile()
	if err != nil {
		log.Fatal(err)
	}
	defer commitsFile.Close()

	branchPathSlices := strings.Split(commitsFile.Name(), "/")
	fmt.Printf("Logs on branch %s\n", branchPathSlices[len(branchPathSlices)-1])

	for _, commit := range fileutils.ReadLines(commitsFile) {
		_, author, message := readCommitObjectContent(commit)
		fmt.Printf("%s    (author %s)    %s\n", commit, author, message)
	}
}

// TODO: refactor this shit
func readCommitObjectContent(hash string) (treeHash string, author string, message string) {
	objectFile, err := os.Open(fmt.Sprintf("%s/%s/%s", fileutils.ObjectsDir, hash[0:2], hash[2:]))
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(objectFile)
	scanner.Split(bufio.ScanLines)
	var contents []string
	for scanner.Scan() {
		contents = append(contents, scanner.Text())
	}
	objectFile.Close()

	treeHash = strings.Split(contents[0], " ")[1]
	author = strings.Split(contents[1], " ")[1]
	message = contents[3]
	return
}

package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aziflaj/gogot/core"
)

// Log ...
func Log(args []string) {
	commitsFile := commitsFile()
	defer commitsFile.Close()

	branchPathSlices := strings.Split(commitsFile.Name(), "/")
	fmt.Printf("Logs on branch %s\n", branchPathSlices[len(branchPathSlices)-1])

	commits := commitsList(commitsFile)
	for _, commit := range commits {
		_, author, message := readCommitObjectContent(commit)
		fmt.Printf("%s    (author %s)    %s\n", commit, author, message)
	}
}

func commitsFile() *os.File {
	currentBranchPath := currentRef()
	indexFile, err := os.Open(fmt.Sprintf("%s/%s", core.GogotDir, currentBranchPath))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return indexFile
}

func commitsList(file *os.File) []string {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var commits []string
	for scanner.Scan() {
		commits = append(commits, scanner.Text())
	}
	return commits
}

func readCommitObjectContent(hash string) (treeHash string, author string, message string) {
	objectFile, err := os.Open(fmt.Sprintf("%s/%s/%s", core.ObjectsDir, hash[0:2], hash[2:]))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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

func readObjectContent(hash string) string {
	objectFile, err := os.Open(fmt.Sprintf("%s/%s/%s", core.ObjectsDir, hash[0:2], hash[2:]))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(objectFile)
	scanner.Split(bufio.ScanLines)
	var contents []string
	for scanner.Scan() {
		contents = append(contents, scanner.Text())
	}
	objectFile.Close()

	return strings.Join(contents, "\n")
}

func readBlobContent(hash string) []byte {
	objectFile, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s", core.ObjectsDir, hash[0:2], hash[2:]))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return objectFile
}

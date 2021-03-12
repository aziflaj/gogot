package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aziflaj/gogot/indextree"
)

// Commit ...
func Commit(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: gogot commit \"message\"")
		os.Exit(1)
	}

	indexTree := buildIndexTree()
	rootSha := buildObjectTree("root", *indexTree)
	commitSha := buildCommit(rootSha, strings.Join(args, " "))
	updateRef(commitSha)
	clearIndex()
}

func buildIndexTree() *indextree.IndexTree {
	indexTree := indextree.New()

	indexFile, err := os.Open(indexPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer indexFile.Close()
	scanner := bufio.NewScanner(indexFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		splitLine := strings.Split(scanner.Text(), " ")
		sha := splitLine[0]
		path := splitLine[1]
		indexTree.AddPath(path, sha)
	}

	return indexTree
}

func buildObjectTree(name string, tree indextree.IndexTree) string {
	sha := hashContent([]byte(time.Now().String() + name))

	objectDirPath := fmt.Sprintf("%s/%s", objectsDir, sha[0:2])
	os.Mkdir(objectDirPath, 0755)
	objectPath := fmt.Sprintf("%s/%s", objectDirPath, sha[2:])
	f, err := os.OpenFile(objectPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	for _, child := range tree.Children {
		if child.Sha != "" { // is a file
			f.WriteString(fmt.Sprintf("blob %s %s\n", child.Sha, child.Name))
		} else { // is a dir
			dirSha := buildObjectTree(child.Name, *child)
			f.WriteString(fmt.Sprintf("tree %s %s\n", dirSha, child.Name))
		}
	}

	return sha
}

func buildCommit(treeSha string, commitMsg string) string {
	whoami := exec.Command("whoami")
	var out bytes.Buffer
	whoami.Stdout = &out
	err := whoami.Run()
	if err != nil {
		fmt.Println("WTF?!")
		os.Exit(1)
	}

	committer := out.String()
	sha := hashContent([]byte(time.Now().String() + committer))

	objectDirPath := fmt.Sprintf("%s/%s", objectsDir, sha[0:2])
	os.Mkdir(objectDirPath, 0755)
	objectPath := fmt.Sprintf("%s/%s", objectDirPath, sha[2:])
	f, err := os.OpenFile(objectPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	f.WriteString(fmt.Sprintf("tree %s\n", treeSha))
	f.WriteString(fmt.Sprintf("author %s\n", committer))
	f.WriteString(fmt.Sprintf("\n%s\n", commitMsg))

	return sha
}

func updateRef(sha string) {
	headFile, err := os.Open(headPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reader := bufio.NewReader(headFile)
	content, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	headFile.Close()

	splitContent := strings.Split(string(content), "/")
	branchName := splitContent[len(splitContent)-1]
	branchPath := fmt.Sprintf("%s/%s", gogotDir, branchName)

	f, err := os.OpenFile(branchPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer f.Close()
	f.WriteString(sha)
}

func clearIndex() error {
	return os.Truncate(indexPath, 0)
}

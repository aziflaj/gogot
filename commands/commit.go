package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/aziflaj/gogot/core"
	"github.com/aziflaj/gogot/files"
	"github.com/aziflaj/gogot/gogot_object"
	"github.com/aziflaj/gogot/index_tree"
)

// Commit ...
func Commit(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: gogot commit <message>")
		os.Exit(1)
	}

	indexTree := buildIndexTree()
	rootHash := buildObjectTree("root", *indexTree)
	commitHash := buildCommitObject(rootHash, strings.Join(args, " "))
	updateRef(commitHash)
	clearIndex()
}

func buildIndexTree() *index_tree.IndexTree {
	indexFile, err := os.Open(files.IndexFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	indexTree := index_tree.BuildFromFile(indexFile)
	indexFile.Close()

	return indexTree
}

func buildObjectTree(name string, tree index_tree.IndexTree) string {
	object, err := gogot_object.CreateFromString(name)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer object.FlushAndClose()

	for _, child := range tree.Children {
		if child.Hash != "" { // is a file
			object.AddBlob(child)
		} else { // is a dir
			dirHash := buildObjectTree(child.Name, *child)
			object.AddTree(child, dirHash)
		}
	}

	return object.Hash
}

func buildCommitObject(treeHash string, commitMsg string) string {
	committer := currentUser()
	commit := core.NewCommitObject(treeHash, committer, commitMsg)
	err := commit.Commit()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return commit.Hash
}

func updateRef(hash string) {
	ref := currentRef()
	branchPath := fmt.Sprintf("%s/%s", files.GogotDir, ref)

	branchFile, err := os.OpenFile(branchPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	branchFile.WriteString(hash + "\n")
	branchFile.Close()
}

func clearIndex() error {
	return os.Truncate(files.IndexFilePath, 0)
}

func currentUser() string {
	whoami := exec.Command("whoami")
	var out bytes.Buffer
	whoami.Stdout = &out
	err := whoami.Run()
	if err != nil {
		fmt.Println("Can't read user name")
		os.Exit(1)
	}

	return strings.Split(out.String(), "\n")[0]
}

func currentRef() string {
	headFile, err := os.Open(files.HeadFilePath)
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

	ref := strings.Split(string(content), ": ")[1]

	return ref
}

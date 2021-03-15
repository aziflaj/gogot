package commands

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/aziflaj/gogot/core"
	"github.com/aziflaj/gogot/fileutils"
)

// Commit ...
func Commit(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: gogot commit <message>")
		os.Exit(1)
	}

	indexTree := buildIndexTree()
	rootHash, err := indexTree.BuildObjectTree("root")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	commitHash := buildCommitObject(rootHash, strings.Join(args, " "))
	updateRef(commitHash)
	clearIndex()
}

func buildIndexTree() *core.IndexTree {
	indexFile, err := os.Open(fileutils.IndexFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	indexTree := core.BuildIndexFromFile(indexFile)
	indexFile.Close()

	return indexTree
}

func buildCommitObject(treeHash string, commitMsg string) string {
	committer := currentUser()
	commit := core.NewCommitObject(treeHash, committer, commitMsg)
	err := commit.Commit()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return commit.ID
}

func updateRef(hash string) {
	ref, err := fileutils.CurrentRef()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	branchPath := fmt.Sprintf("%s/%s", fileutils.GogotDir, ref)
	branchFile, err := os.OpenFile(branchPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	branchFile.WriteString(hash + "\n")
	branchFile.Close()
}

func clearIndex() error {
	return os.Truncate(fileutils.IndexFilePath, 0)
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

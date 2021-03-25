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
		log.Fatal("Usage: gogot commit <message>")
	}

	indexTree := buildIndexTree()
	rootHash, err := indexTree.BuildObjectTree("root")
	if err != nil {
		log.Fatal(err)
	}

	commitHash := buildCommitObject(rootHash, strings.Join(args, " "))
	updateRef(commitHash)
	clearIndex()
}

func buildIndexTree() *core.IndexTree {
	indexFile, err := os.Open(fileutils.IndexFilePath)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	return commit.ID
}

func updateRef(hash string) {
	ref, err := fileutils.CurrentRef()
	if err != nil {
		log.Fatal(err)
	}

	branchPath := fmt.Sprintf("%s/%s", fileutils.GogotDir, ref)
	branchFile, err := os.OpenFile(branchPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal("Can't read user name")
	}

	return strings.Split(out.String(), "\n")[0]
}

package commands

import (
	"fmt"
	"log"

	"github.com/aziflaj/gogot/core"
	"github.com/aziflaj/gogot/fileutils"
)

// TimeMachine ...
func TimeMachine(args []string) {
	if len(args) < 2 {
		log.Fatal("Usage: gogot time-machine <commit-id> <file-path>")
	}

	commitID, filePath := args[0], args[1]
	commit, err := core.FindCommitWithID(commitID)
	if err != nil {
		log.Fatal(err)
	}

	pastIndexTree, err := core.BuildIndexFromCommit(commit.TreeHash, ".")
	if err != nil {
		log.Fatal(err)
	}

	child := pastIndexTree.FindChildByPath(filePath)
	if child == nil {
		log.Fatalf("File %s not found\n", filePath)
	}
	blob, err := fileutils.ReadBlobContents(child.Hash)
	if err != nil {
		log.Fatal(err)
	}

	fileContent, err := core.DecompressBytes(blob)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fileContent)
}

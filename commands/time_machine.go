package commands

import (
	"fmt"
	"os"

	"github.com/aziflaj/gogot/core"
	"github.com/aziflaj/gogot/fileutils"
)

// TimeMachine ...
func TimeMachine(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: gogot time-machine <commit-id> <file-path>")
		os.Exit(1)
	}

	commitID, filePath := args[0], args[1]
	commit, err := core.FindCommitWithID(commitID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pastIndexTree, err := core.BuildIndexFromCommit(commit.TreeHash, ".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	child := pastIndexTree.FindChildByPath(filePath)
	if child == nil {
		fmt.Printf("File %s not found\n", filePath)
		os.Exit(0)
	}
	blob, err := fileutils.ReadBlobContents(child.Hash)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fileContent, err := core.DecompressBytes(blob)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(fileContent)
}

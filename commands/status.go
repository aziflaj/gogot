package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aziflaj/gogot/core"
	"github.com/aziflaj/gogot/fileutils"
)

// Status ...
// TODO: Add the following as feature requests:
// gogot rm <file>          - opposite of gogot add
// gogot rollback <file>    - persistent time-machine (last commit)
func Status(args []string) {
	branch, err := currentBranch()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("On branch %s\n", branch)

	printIndexedChanges()
	fmt.Print("\n")

	commits, err := commitsInBranch()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	paths, err := fileutils.AllPaths(".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var trackedFiles, untrackedFiles []string

nextPath:
	for _, filePath := range paths {
		for _, commit := range commits {
			commitTree, err := core.BuildIndexFromCommit(commit.TreeHash, ".")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			child := commitTree.FindChildByPath(filePath)
			// fmt.Println("child")
			// fmt.Println(child)
			// fmt.Println("commitTree")
			// fmt.Println(commitTree)
			// fmt.Println("filePath")
			// fmt.Println(filePath)

			if child != nil {
				// TODO: check if current blob matches child's blob
				trackedFiles = append(trackedFiles, filePath)
				continue nextPath
			}
		}
		untrackedFiles = append(untrackedFiles, filePath)
	}

	// tracked files, unindexed (prev commit, not in index)
	fmt.Println("Files not added to index:")
	fmt.Println("    (use \"gogot add/rm <file>\") to update what will be committed")
	fmt.Println("    (use \"gogot rollback <file>\") to unstage")
	for _, file := range trackedFiles {
		fmt.Printf("\t%s\n", file)
	}

	fmt.Print("\n")

	// untracked files (they shouldn't be in the prev commit)
	fmt.Println("Untracked files:")
	fmt.Println("    (use \"gogot add <file>\") to include in the commit")
	for _, file := range untrackedFiles {
		fmt.Printf("\t%s\n", file)
	}

	fmt.Println("nothing to commit, working tree clean")
}

func printIndexedChanges() {
	stagedFiles, err := filesInIndex()
	if err != nil {
		fmt.Println("No commits yet")
		os.Exit(0)
	}

	if len(stagedFiles) > 0 {
		fmt.Println("Files to be committed:")
		fmt.Println("    (use \"gogot rollback <file>\") to unstage")
		for _, file := range stagedFiles {
			fmt.Printf("\t%s\n", file)
		}
	}
}

func commitsInBranch() (commits []*core.CommitObject, err error) {
	commitsFile, err := fileutils.CurrentBranchCommitsFile()
	if err != nil {
		return nil, err
	}

	commitIds := fileutils.ReadLines(commitsFile)
	commitsCount := len(commitIds)
	commits = make([]*core.CommitObject, commitsCount)
	for i := 0; i < commitsCount; i++ {
		commits[i], err = core.FindCommitWithID(commitIds[commitsCount-i-1])
	}

	return
}

func currentBranch() (string, error) {
	ref, err := fileutils.CurrentRef()
	if err != nil {
		return "", err
	}

	splitRef := strings.Split(ref, "/")

	return splitRef[len(splitRef)-1], nil
}

func filesInIndex() (files []string, err error) {
	indexFile, err := os.Open(fileutils.IndexFilePath)
	if err != nil {
		return []string{}, err
	}

	for _, index := range fileutils.ReadLines(indexFile) {
		hashAndPath := strings.Split(index, " ")
		files = append(files, hashAndPath[1])
	}

	return files, nil
}

// true if files are different
func compareFiles(filePath string, child *core.IndexTree) bool {
	currentContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hash := core.HashBytes(currentContent)
	return hash == child.Hash
}

func isFileIndexed(filePath string, commit *core.CommitObject) bool {
	parentCommit, err := commit.Parent()
	if err != nil || parentCommit == nil {
		return false
	}

	commitTree, err := core.BuildIndexFromCommit(parentCommit.TreeHash, ".")
	if err != nil {
		return false
	}

	child := commitTree.FindChildByPath(filePath)
	if child == nil {
		return isFileIndexed(filePath, parentCommit)
	}

	return true
}

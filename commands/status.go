package commands

import (
	"fmt"
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

	prevCommit, err := previousCommitTree()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	paths, err := fileutils.AllPaths(".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, filePath := range paths {
		child := prevCommit.FindChildByPath(filePath)
		if child != nil {
			// check if changed after indexed
			fmt.Println("tracked " + filePath)
		} else {
			// fmt.Println("untracked " + filePath)
			// untracked = append(untracked, file.Name())
		}
	}

	// untrackedFiles, unindexedFiles, err := categorizeFiles(allFiles, prevCommit)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// tracked files, unindexed (prev commit, not in index)
	fmt.Println("Files not added to index:")
	fmt.Println("    (use \"gogot add/rm <file>\") to update what will be committed")
	fmt.Println("    (use \"gogot rollback <file>\") to unstage")
	// for _, file := range unindexedFiles {
	// 	fmt.Printf("\t%s\n", file)
	// }

	fmt.Print("\n")

	// untracked files (they shouldn't be in the prev commit)
	fmt.Println("Untracked files:")
	fmt.Println("    (use \"gogot add <file>\") to include in the commit")
	// for _, file := range untrackedFiles {
	// 	fmt.Printf("\t%s\n", file)
	// }

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

func previousCommitTree() (*core.IndexTree, error) {
	commitsFile, err := fileutils.CurrentBranchCommitsFile()
	if err != nil {
		return nil, err
	}

	commits := fileutils.ReadLines(commitsFile)
	prevCommitId := commits[len(commits)-1]
	commit, err := core.FindCommitWithID(prevCommitId)
	if err != nil {
		return nil, err
	}

	return core.BuildIndexFromCommit(commit.TreeHash, ".")
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

// func categorizeFiles(allFiles []*os.File, commitTree *core.IndexTree) ([]string, []string, error) {
// 	var untracked, unindexed []string

// categorize:
// 	for _, file := range allFiles {
// 		if file.Name() == fileutils.GogotDir {
// 			continue
// 		}

// 		for _, pattern := range patterns {
// 			if match, _ := path.Match(pattern, file.Name()); match {
// 				continue categorize
// 			}
// 		}

// 		info, err := os.Stat(file.Name())
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}

// 		if info.IsDir() {
// 			// handle dir

// 			// subfiles, _ := ioutil.ReadDir(file.Name())
// 			// for _, subfile := range subfiles {

// 			// }
// 		} else {
// 			// handle File
// 			child := commitTree.FindChildByPath(file.Name())
// 			if child != nil {
// 				untracked = append(untracked, file.Name())
// 			} else {
// 				// check if changed after indexed
// 			}
// 		}
// 	}

// 	return untracked, unindexed, nil
// }

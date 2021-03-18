package core

import (
	"fmt"
	"strings"

	"github.com/aziflaj/gogot/fileutils"
)

type CommitObject struct {
	ID       string
	TreeHash string
	Author   string
	Message  string
}

func NewCommitObject(treeHash string, author string, commitMsg string) *CommitObject {
	id := TimedHash(author)

	return &CommitObject{
		ID:       id,
		TreeHash: treeHash,
		Author:   author,
		Message:  commitMsg,
	}
}

func FindCommitWithID(id string) (*CommitObject, error) {
	content, err := fileutils.ReadCommitContents(id)
	if err != nil {
		return nil, err
	}

	splitContent := strings.Split(content, "\n")
	return &CommitObject{
		ID:       id,
		TreeHash: strings.Split(splitContent[0], "tree ")[1],
		Author:   strings.Split(splitContent[1], "author ")[1],
		Message:  splitContent[3],
	}, nil
}

func (obj *CommitObject) Commit() error {
	file, err := fileutils.CreateAndOpenCommitFile(obj.ID)
	if err != nil {
		return err
	}

	defer file.Close()
	file.WriteString(fmt.Sprintf("tree %s\n", obj.TreeHash))
	file.WriteString(fmt.Sprintf("author %s\n", obj.Author))
	file.WriteString(fmt.Sprintf("\n%s\n", obj.Message))

	return nil
}

func (obj *CommitObject) Parent() (*CommitObject, error) {
	commitsFile, err := fileutils.CurrentBranchCommitsFile()
	if err != nil {
		return nil, err
	}

	commits := fileutils.ReadLines(commitsFile)
	for index, commitId := range commits {
		if commitId == obj.ID && index > 0 {
			return FindCommitWithID(commits[index-1])
		}
	}

	return nil, nil
}

func (obj *CommitObject) String() string {
	return fmt.Sprintf("[Commit: %s]", obj.ID)
}

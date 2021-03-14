package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/aziflaj/gogot/fileutils"
)

type CommitObject struct {
	Hash     string
	TreeHash string
	Author   string
	Message  string
}

func NewCommitObject(treeHash string, author string, commitMsg string) *CommitObject {
	hash := HashBytes([]byte(time.Now().String() + author))

	return &CommitObject{
		Hash:     hash,
		TreeHash: treeHash,
		Author:   author,
		Message:  commitMsg,
	}
}

func CommitObjectFromHash(hash string) (*CommitObject, error) {
	content, err := fileutils.ReadCommitContents(hash)
	if err != nil {
		return nil, err
	}

	splitContent := strings.Split(content, "\n")
	return &CommitObject{
		Hash:     hash,
		TreeHash: strings.Split(splitContent[0], "tree ")[1],
		Author:   strings.Split(splitContent[1], "author ")[1],
		Message:  splitContent[3],
	}, nil
}

func (obj *CommitObject) Commit() error {
	file, err := fileutils.CreateAndOpenCommitFile(obj.Hash)
	if err != nil {
		return err
	}

	defer file.Close()
	file.WriteString(fmt.Sprintf("tree %s\n", obj.TreeHash))
	file.WriteString(fmt.Sprintf("author %s\n", obj.Author))
	file.WriteString(fmt.Sprintf("\n%s\n", obj.Message))

	return nil
}

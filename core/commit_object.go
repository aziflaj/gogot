package core

import (
	"fmt"
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

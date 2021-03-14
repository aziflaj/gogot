package core

import (
	"fmt"
	"os"
	"time"

	"github.com/aziflaj/gogot/fileutils"
)

type GogotObject struct {
	Hash       string
	ObjectFile os.File
}

func CreateObjectFromString(str string) (*GogotObject, error) {
	hash := HashBytes([]byte(time.Now().String() + str))
	file, err := fileutils.CreateAndOpenCommitFile(hash)
	if err != nil {
		return nil, err
	}

	return &GogotObject{Hash: hash, ObjectFile: *file}, nil
}

func (obj *GogotObject) AddBlob(blob *IndexTree) {
	obj.ObjectFile.WriteString(fmt.Sprintf("blob %s %s\n", blob.Hash, blob.Name))
}

func (obj *GogotObject) AddTree(tree *IndexTree, hash string) {
	obj.ObjectFile.WriteString(fmt.Sprintf("tree %s %s\n", hash, tree.Name))
}

func (obj *GogotObject) FlushAndClose() {
	obj.ObjectFile.Close()
}

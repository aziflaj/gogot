package gogot_object

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/aziflaj/gogot/core"
	"github.com/aziflaj/gogot/index_tree"
)

type GogotObject struct {
	Hash       string
	ObjectFile os.File
}

func CreateFromString(str string) (*GogotObject, error) {
	hash := hashContent([]byte(time.Now().String() + str))
	file, err := createAndOpenFile(hash)
	if err != nil {
		return nil, err
	}

	return &GogotObject{Hash: hash, ObjectFile: *file}, nil
}

func (obj *GogotObject) Write(str string) {
	obj.ObjectFile.WriteString(str)
}

func (obj *GogotObject) AddBlob(blob *index_tree.IndexTree) {
	obj.ObjectFile.WriteString(fmt.Sprintf("blob %s %s\n", blob.Hash, blob.Name))
}

func (obj *GogotObject) AddTree(tree *index_tree.IndexTree, hash string) {
	obj.ObjectFile.WriteString(fmt.Sprintf("tree %s %s\n", hash, tree.Name))
}

func (obj *GogotObject) FlushAndClose() {
	obj.ObjectFile.Close()
}

func createAndOpenFile(hash string) (file *os.File, err error) {
	objectDirPath := fmt.Sprintf("%s/%s", core.ObjectsDir, hash[0:2])
	os.Mkdir(objectDirPath, 0755)
	objectPath := fmt.Sprintf("%s/%s", objectDirPath, hash[2:])
	file, err = os.OpenFile(objectPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return
}

func hashContent(content []byte) string {
	hasher := sha1.New()
	hasher.Write(content)
	sha1Bytes := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return string(sha1Bytes)
}

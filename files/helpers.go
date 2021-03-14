package files

import (
	"fmt"
	"os"
)

func CreateAndOpenCommitFile(hash string) (file *os.File, err error) {
	objectDirPath := fmt.Sprintf("%s/%s", ObjectsDir, hash[0:2])
	os.Mkdir(objectDirPath, 0755)
	objectPath := fmt.Sprintf("%s/%s", objectDirPath, hash[2:])
	file, err = os.OpenFile(objectPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return
}

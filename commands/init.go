package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/aziflaj/gogot/fileutils"
)

// Init ...
func Init(args []string) {
	if len(args) < 1 {
		log.Fatal("Usage: gogot init <path>")
	}

	repoName := args[0]

	fmt.Println("Initalizing new Gogot repo")

	baseRepoPath := fmt.Sprintf("%s/%s", repoName, fileutils.GogotDir)
	os.MkdirAll(baseRepoPath, 0755)
	createObjectsDir(baseRepoPath)
	createRefsDir(baseRepoPath)
	initializeHead(baseRepoPath)

	log.Printf("Gogot repo initialized in %s\n", baseRepoPath)
}

func createObjectsDir(basePath string) {
	objectsDir := fmt.Sprintf("%s/objects", basePath) // Should use this: fileutils.ObjectsDir
	infoDir := fmt.Sprintf("%s/info", objectsDir)
	packDir := fmt.Sprintf("%s/pack", objectsDir)

	err := os.Mkdir(objectsDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(infoDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(packDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func createRefsDir(basePath string) {
	refsDir := fmt.Sprintf("%s/refs", basePath) // fileutils.RefsDir
	headsDir := fmt.Sprintf("%s/heads", refsDir)
	tagsDir := fmt.Sprintf("%s/tags", refsDir)

	os.MkdirAll(refsDir, 0755)
	os.MkdirAll(headsDir, 0755)
	os.MkdirAll(tagsDir, 0755)
}

func initializeHead(basePath string) {
	headFilePath := fmt.Sprintf("%s/HEAD", basePath)
	file, err := os.Create(headFilePath) // fileutils.HeadFilePath
	if err != nil {
		cleanup(basePath)
	}
	defer file.Close()

	file.WriteString("ref: refs/heads/main")
}

func cleanup(path string) {
	// os.RemoveAll(fileutils.GogotDir)
	os.RemoveAll(path)
	log.Fatal("Something went wrong")
}

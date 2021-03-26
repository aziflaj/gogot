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
	createObjectsDir()
	createRefsDir()
	initializeHead()

	log.Printf("Gogot repo initialized in %s\n", baseRepoPath)
}

func createObjectsDir() {
	infoDir := fmt.Sprintf("%s/info", fileutils.ObjectsDir)
	packDir := fmt.Sprintf("%s/pack", fileutils.ObjectsDir)

	os.MkdirAll(fileutils.ObjectsDir, 0755)
	os.MkdirAll(infoDir, 0755)
	os.MkdirAll(packDir, 0755)
}

func createRefsDir() {
	headsDir := fmt.Sprintf("%s/heads", fileutils.RefsDir)
	tagsDir := fmt.Sprintf("%s/tags", fileutils.RefsDir)

	os.MkdirAll(fileutils.RefsDir, 0755)
	os.MkdirAll(headsDir, 0755)
	os.MkdirAll(tagsDir, 0755)
}

func initializeHead() {
	file, err := os.Create(fileutils.HeadFilePath)
	if err != nil {
		cleanup()
	}
	defer file.Close()

	file.WriteString("ref: refs/heads/main")
}

func cleanup() {
	os.RemoveAll(fileutils.GogotDir)
	log.Fatal("Something went wrong")
}

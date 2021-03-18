package commands

import (
	"fmt"
	"os"

	"github.com/aziflaj/gogot/fileutils"
)

// Init ...
func Init(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: gogot init <path>")
		os.Exit(1)
	}

	repoName := args[0]

	fmt.Println("Initalizing new Gogot repo")

	baseRepoPath := fmt.Sprintf("%s/%s", repoName, fileutils.GogotDir)
	err := os.MkdirAll(baseRepoPath, 0755)
	if err != nil {
		fmt.Print(err)
		fmt.Println("Repo already exists")
		os.Exit(1)
	}

	createObjectsDir()
	createRefsDir()
	initializeHead()

	fmt.Printf("Gogot repo initialized in %s\n", baseRepoPath)
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
	fmt.Println("Something went wrong")
	os.RemoveAll(fileutils.GogotDir)
	os.Exit(1)
}

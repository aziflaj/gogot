package commands

import (
	"fmt"
	"os"

	"github.com/aziflaj/gogot/files"
)

// Init ...
func Init(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: gogot init <path>")
		os.Exit(1)
	}

	repoName := args[0]

	fmt.Println("Initalizing new Gogot repo")

	baseRepoPath := fmt.Sprintf("%s/%s", repoName, files.GogotDir)
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

func createObjectsDir() error {
	infoDir := fmt.Sprintf("%s/info", files.ObjectsDir)
	packDir := fmt.Sprintf("%s/pack", files.ObjectsDir)

	var err error

	err = os.MkdirAll(files.ObjectsDir, 0755)
	err = os.MkdirAll(infoDir, 0755)
	err = os.MkdirAll(packDir, 0755)

	return err
}

func createRefsDir() error {
	headsDir := fmt.Sprintf("%s/heads", files.RefsDir)
	tagsDir := fmt.Sprintf("%s/tags", files.RefsDir)

	var err error

	err = os.MkdirAll(files.RefsDir, 0755)
	err = os.MkdirAll(headsDir, 0755)
	err = os.MkdirAll(tagsDir, 0755)

	return err
}

func initializeHead() {
	file, err := os.Create(files.HeadFilePath)
	if err != nil {
		cleanup()
	}
	defer file.Close()

	file.WriteString("ref: refs/heads/main")
}

func cleanup() {
	fmt.Println("Something went wrong")
	os.RemoveAll(files.GogotDir)
	os.Exit(1)
}

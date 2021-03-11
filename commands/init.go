package commands

import (
	"fmt"
	"os"
)

// Init ...
func Init(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: gogot init [PATH]")
		os.Exit(1)
	}

	repoName := args[0]

	fmt.Println("Initalizing new Gogot repo")

	baseRepoPath := fmt.Sprintf("%s/.gogot", repoName)
	err := os.MkdirAll(baseRepoPath, 0755)
	if err != nil {
		fmt.Print(err)
		fmt.Println("Repo already exists")
		os.Exit(1)
	}

	createObjectsDir(baseRepoPath)
	createRefsDir(baseRepoPath)
	initializeHead(baseRepoPath)

	fmt.Printf("Gogot repo initialized in %s\n", baseRepoPath)
}

func createObjectsDir(path string) error {
	objectsDir := fmt.Sprintf("%s/objects", path)
	infoDir := fmt.Sprintf("%s/info", objectsDir)
	packDir := fmt.Sprintf("%s/pack", objectsDir)

	var err error

	err = os.MkdirAll(objectsDir, 0755)
	err = os.MkdirAll(infoDir, 0755)
	err = os.MkdirAll(packDir, 0755)

	return err
}

func createRefsDir(path string) error {
	refsDir := fmt.Sprintf("%s/refs", path)
	headsDir := fmt.Sprintf("%s/heads", refsDir)
	tagsDir := fmt.Sprintf("%s/tags", refsDir)

	var err error

	err = os.MkdirAll(refsDir, 0755)
	err = os.MkdirAll(headsDir, 0755)
	err = os.MkdirAll(tagsDir, 0755)

	return err
}

func initializeHead(path string) {
	filename := fmt.Sprintf("%s/HEAD", path)
	file, err := os.Create(filename)

	if err != nil {
		cleanup(path)
	}
	defer file.Close()

	file.WriteString("ref: refs/heads/master")
}

func cleanup(path string) {
	fmt.Println("Something went wrong")
	os.RemoveAll(path)
	os.Exit(1)
}

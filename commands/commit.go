package commands

import (
	"fmt"
	"os"
)

// Commit ...
func Commit(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: gogot commit -m \"message\"")
		os.Exit(1)
	}
}

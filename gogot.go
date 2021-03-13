package main

import (
	"fmt"
	"os"

	"github.com/aziflaj/gogot/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gogot <command> [<args>]")
		os.Exit(1)
	}

	command, args := os.Args[1], os.Args[2:]

	switch command {
	case "init":
		commands.Init(args)
	case "add":
		commands.Add(args)
	case "commit":
		commands.Commit(args)
	case "log":
		commands.Log(args)
	default:
		fmt.Printf("Unknown command: %v\n", command)
		os.Exit(1)
	}
}

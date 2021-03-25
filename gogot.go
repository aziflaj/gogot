package main

import (
	"log"
	"os"

	"github.com/aziflaj/gogot/commands"
	"github.com/aziflaj/gogot/fileutils"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		log.Fatal("Usage: gogot <command> [<args>]")
	}

	command, args := os.Args[1], os.Args[2:]

	if _, err := os.Stat(fileutils.GogotDir); os.IsNotExist(err) && command != "init" {
		log.Fatal("Not a Gogot repository")
	}

	switch command {
	case "init":
		commands.Init(args)
	case "add":
		commands.Add(args)
	case "commit":
		commands.Commit(args)
	case "log":
		commands.Log(args)
	case "time-machine":
		commands.TimeMachine(args)
	case "status":
		commands.Status(args)
	default:
		log.Fatalf("Unknown command: %v\n", command)
	}
}

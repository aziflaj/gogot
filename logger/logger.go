package logger

import (
	"log"
	"os"
)

func SetupLogger() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
}

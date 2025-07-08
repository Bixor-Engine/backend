package main

import (
	"os"

	"github.com/Bixor-Engine/backend/cmd/bixor/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}

package main

import (
	"os"

	"Edot/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(0)
	}
}

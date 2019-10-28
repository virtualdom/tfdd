package main

import (
	"os"

	"github.com/virtualdom/tfdd/pkg/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}


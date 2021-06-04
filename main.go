package main

import (
	"os"

	"github.com/devops-kung-fu/lucha/cmd"
)

func main() {
	defer os.Exit(0)
	cmd.Execute()
}

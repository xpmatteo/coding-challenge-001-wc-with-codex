package main

import (
	"os"

	"wc/internal/cli"
)

func main() {
	app := cli.App{}
	exitCode := app.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

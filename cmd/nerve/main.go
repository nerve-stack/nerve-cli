package main

import (
	"os"

	"github.com/nerve-stack/nerve-cli/internal/cmd"
)

func main() {
	os.Exit(cmd.Do(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

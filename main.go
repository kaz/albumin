package main

import (
	"os"

	"github.com/kaz/albumin/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		os.Exit(1)
	}
}

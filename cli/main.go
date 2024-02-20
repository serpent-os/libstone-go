package main

import (
	"fmt"
	"os"

	"github.com/serpent-os/libstone/internal/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

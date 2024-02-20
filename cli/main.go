// SPDX-FileCopyrightText: 2024 Serpent OS Developers
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/serpent-os/libstone-go/internal/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

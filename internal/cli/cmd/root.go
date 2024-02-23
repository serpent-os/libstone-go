// SPDX-FileCopyrightText: 2024 Serpent OS Developers
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"github.com/alecthomas/kong"
)

// Version is the version of the library and the command line interface.
var Version string

type globalFlags struct {
	Version kong.VersionFlag `help:"Prints version and exits."`
}

type cli struct {
	globalFlags

	Inspect cmdInspect `cmd:"" help:"Inspect stone package contents."`
}

// Run runs the command line interface.
func Run() error {
	var cli cli
	ctx := kong.Parse(
		&cli,
		kong.Name("libstone"),
		kong.Description("A Golang implementation for stone binary packages"),
		kong.Vars{"version": Version},
	)
	return ctx.Run(&cli.globalFlags)
}

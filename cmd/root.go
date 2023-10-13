package cmd

import (
	"github.com/rebuy-de/rebuy-go-sdk/v6/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	return cmdutil.New(
		"libstone", "A golang implementation for stone binary packages",
		cmdutil.WithLogVerboseFlag(),
		cmdutil.WithVersionCommand(),

		cmdutil.WithSubCommand(cmdutil.New(
			"inspect", "Inspect stone package contents",
			cmdutil.WithRun(Inspect),
		)),
	)
}

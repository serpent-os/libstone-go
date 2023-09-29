//go:build tools

package main

// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
import (
	_ "golang.org/x/tools/cmd/stringer"
)

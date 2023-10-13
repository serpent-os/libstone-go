//go:build tools

package main

// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
import (
	_ "github.com/rebuy-de/rebuy-go-sdk/v6/cmd/buildutil"
	_ "github.com/rebuy-de/rebuy-go-sdk/v6/cmd/cdnmirror"
	_ "golang.org/x/tools/cmd/stringer"
)

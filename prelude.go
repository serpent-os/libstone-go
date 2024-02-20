// SPDX-FileCopyrightText: 2024 Serpent OS Developers
// SPDX-License-Identifier: MPL-2.0

package libstone

import (
	"errors"
	"io"

	"github.com/serpent-os/libstone-go/internal/readers"
)

// Version is the stone format version contained inside the [Prelude].
type Version uint32

const (
	// V1 is the first version of the stone format.
	V1 Version = iota + 1
)

var (
	// ErrNoStone is returned when the magic number doesn't match
	// [MagicNumber].
	ErrNoStone = errors.New("data is not a stone archive")
)

const (
	preludeLen = 32
)

var (
	// magicNumber is the magic number of a stone archive.
	magicNumber = [4]byte{0, 'm', 'o', 's'}
)

// PreludeData is an agnostic array of bytes extending the base Prelude.
// Its meaning varies according to Version.
type PreludeData [24]byte

// Prelude is the header of the stone format.
type Prelude struct {
	Data PreludeData

	// Version is the version of this stone archive.
	Version Version
}

func ReadPrelude(src io.Reader) (Prelude, error) {
	var rawPrelude [preludeLen]byte
	_, err := io.ReadFull(src, rawPrelude[:])
	if err != nil {
		return Prelude{}, err
	}

	wlk := readers.ByteWalker(rawPrelude[:])
	if [4]byte(wlk.Ahead(len(magicNumber))) != magicNumber {
		return Prelude{}, ErrNoStone
	}
	var out Prelude
	out.Data = PreludeData(wlk.Ahead(len(out.Data)))
	out.Version = Version(wlk.Uint32())
	return out, nil
}

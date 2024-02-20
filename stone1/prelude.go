// SPDX-FileCopyrightText: 2024 Serpent OS Developers
// SPDX-License-Identifier: MPL-2.0

package stone1

import (
	"errors"

	"github.com/serpent-os/libstone-go"
	"github.com/serpent-os/libstone-go/internal/readers"
)

var (
	integrityCheck = [21]byte{0, 0, 1, 0, 0, 2, 0, 0, 3, 0, 0, 4, 0, 0, 5, 0, 0, 6, 0, 0, 7}
)

type StoneType uint8

const (
	BinaryStone StoneType = iota + 1
	DeltaStone
	RepositoryStone
	BuildManifestStone
)

type Prelude struct {
	NumPayloads uint16
	StoneType   StoneType
}

func NewPrelude(genericPre libstone.Prelude) (Prelude, error) {
	if genericPre.Version != libstone.V1 {
		return Prelude{}, errors.New("prelude version is not 1")
	}
	if len(genericPre.Data) < len(libstone.PreludeData{}) {
		return Prelude{}, errors.New("insufficient number of bytes to parse a V1 prelude")
	}

	wlk := readers.ByteWalker(genericPre.Data[:])
	var pre Prelude
	pre.NumPayloads = wlk.Uint16()
	if [21]byte(wlk.Ahead(len(integrityCheck))) != integrityCheck {
		return Prelude{}, errors.New("V1 integrity check failed")
	}
	pre.StoneType = StoneType(wlk.Uint8())
	return pre, nil
}

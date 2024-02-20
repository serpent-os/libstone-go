// SPDX-FileCopyrightText: 2024 Serpent OS Developers
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"encoding/hex"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/serpent-os/libstone"
	"github.com/serpent-os/libstone/stone1"
)

type cmdInspect struct {
	Archive string `arg:"" help:"Path of the .stone archive."`
}

func (cmd cmdInspect) Run(globals *globalFlags) error {
	arch, err := os.Open(cmd.Archive)
	if err != nil {
		return err
	}
	defer arch.Close()
	genericPrelude, err := libstone.ReadPrelude(arch)
	if err != nil {
		return err
	}
	prelude, err := stone1.NewPrelude(genericPrelude)
	if err != nil {
		return err
	}
	cache, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}
	defer os.Remove(cache.Name())
	reader := stone1.NewReader(prelude, arch, cache)
	return printArchive(reader)
}

func printArchive(rdr *stone1.Reader) error {
	for rdr.NextPayload() {
		if rdr.Header.Kind != stone1.Meta && rdr.Header.Kind != stone1.Layout {
			fmt.Printf("Inspection of %q record not implemented\n", rdr.Header.Kind)
			continue
		}
		for rdr.NextRecord() {
			switch cast := rdr.Record.(type) {
			case *stone1.MetaRecord:
				printMeta(cast)
			case *stone1.LayoutRecord:
				printLayout(cast)
			}
		}
		if rdr.Err != nil {
			return rdr.Err
		}
	}
	return rdr.Err
}

func printMeta(rec *stone1.MetaRecord) {
	fmt.Printf("%s:\t%s\n", rec.Tag, rec.Field)
}

func printLayout(rec *stone1.LayoutRecord) {
	files := [2]string{string(rec.Entry.Source()), string(rec.Entry.Target())}
	for i := range files {
		if utf8.ValidString(files[i]) {
			continue
		}
		files[i] = hex.EncodeToString([]byte(files[i]))
	}
	fmt.Printf("%s\t-> %s\t[%s]\n", files[0], files[1], rec.Entry.FileType)
}

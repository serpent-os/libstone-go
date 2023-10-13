package header

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

//go:generate -command stringer go run golang.org/x/tools/cmd/stringer
//go:generate stringer -type FileType -output generated_const_names.go

type AgnosticHeader struct {
	/// 4-bytes, BE (u32): Magic to quickly identify a stone file
	Magic [4]byte

	/// 24 bytes, version specific
	//Data [24]byte
	Data V1Data

	/// 4-bytes, BE (u32): Format version used in the container
	Version uint32
}

type V1Data struct {
	NumPayloads    uint16
	IntegrityCheck [21]byte
	FileType       FileType
}

func getStoneMagic() [4]byte {
	return [4]byte{'\x00', 'm', 'o', 's'}
}

func getIntegrityCheck() [21]byte {
	return [21]byte{0, 0, 1, 0, 0, 2, 0, 0, 3, 0, 0, 4, 0, 0, 5, 0, 0, 6, 0, 0, 7}
}

type FileType uint8

const (
	FileTypeUnknown FileType = iota
	FileTypeBinary
	FileTypeDelta
	FileTypeRepository
	FileTypeBuildManifest
)

type Version uint32

func ReadHeader(r io.Reader) (AgnosticHeader, error) {
	agnosticHeader := AgnosticHeader{}

	err := binary.Read(r, binary.BigEndian, &agnosticHeader)
	if err != nil {
		return AgnosticHeader{}, errors.Wrap(err, "Couldn't read header into struct")
	}

	stoneMagic := getStoneMagic()
	integrityCheck := getIntegrityCheck()

	if string(agnosticHeader.Magic[:]) != string(stoneMagic[:]) {
		return AgnosticHeader{}, errors.New("File is no .stone file")
	}

	if string(agnosticHeader.Data.IntegrityCheck[:]) != string(integrityCheck[:]) {
		return AgnosticHeader{}, errors.New("Integrity Check sequence doesn't match")
	}

	if agnosticHeader.Data.FileType > 4 {
		return AgnosticHeader{}, errors.New("Unsupported FileType")
	}

	return agnosticHeader, nil
}

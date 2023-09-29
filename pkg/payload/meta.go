package payload

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/sirupsen/logrus"
)

type Meta struct {
	Tag  RecordTag
	Type RecordType
}

type Dependency uint8

//go:generate -command stringer go run golang.org/x/tools/cmd/stringer
//go:generate stringer -type RecordType,RecordTag
const (
	PackageName Dependency = iota
	SharedLibrary
	PkgConfig
	Interpreter
	CMake
	Python
	Binary
	SystemBinary
	PkgConfig32
)

type RecordTag uint16

const (
	// Name of the package
	RecordTagName RecordTag = iota + 1
	// Architecture of the package
	RecordTagArchitecture
	// Version of the package
	RecordTagVersion
	// Summary of the package
	RecordTagSummary
	// Description of the package
	RecordTagDescription
	// Homepage for the package
	RecordTagHomepage
	// ID for the source package, used for grouping
	RecordTagSourceID
	// Runtime dependencies
	RecordTagDepends
	// Provides some capability or name
	RecordTagProvides
	// Conflicts with some capability or name
	RecordTagConflicts
	// Release number for the package
	RecordTagRelease
	// SPDX license identifier
	RecordTagLicense
	// Currently recorded build number
	RecordTagBuildRelease
	// Repository index specific (relative URI)
	RecordTagPackageURI
	// Repository index specific (Package hash)
	RecordTagPackageHash
	// Repository index specific (size on disk)
	RecordTagPackageSize
	// A Build Dependency
	RecordTagBuildDepends
	// Upstream URI for the source
	RecordTagSourceURI
	// Relative path for the source within the upstream URI
	RecordTagSourcePath
	// Ref/commit of the upstream source
	RecordTagSourceRef
)

type RecordType uint8

const (
	RecordTypeUnknown RecordType = iota
	RecordTypeInt8
	RecordTypeUint8
	RecordTypeInt16
	RecordTypeUint16
	RecordTypeInt32
	RecordTypeUint32
	RecordTypeInt64
	RecordTypeUint64
	RecordTypeString
	RecordTypeDependency
	RecordTypeProvider
)

type MetaRecord struct {
	Length     uint32
	RecordTag  RecordTag
	RecordType RecordType
	Padding    byte
}

func ReadData[T any](input io.Reader) (T, error) {
	var output T
	err := binary.Read(input, binary.BigEndian, &output)
	if err != nil {
		return output, err
	}
	return output, nil
}

func DecodeMetaPayload(payload []byte, records int) error {
	reader := bytes.NewBuffer(payload)
	for i := 0; i < records; i++ {
		record := MetaRecord{}

		err := binary.Read(reader, binary.BigEndian, &record)
		if err != nil {
			return err
		}

		data, err := switchstuff(reader, record.RecordType)
		if err != nil {
			return err
		}

		logrus.Printf("Payload %d, Length %d, Tag %s, Type %s", i, record.Length, record.RecordTag.String(), record.RecordType.String())
		logrus.Printf("%v", data)
	}
	return nil
}

func switchstuff(buf *bytes.Buffer, recordType RecordType) (any, error) {
	switch recordType {
	case RecordTypeString:
		return buf.ReadString('\x00')
	case RecordTypeUint64:
		return ReadData[uint64](buf)
	default:
		return nil, errors.New("Idiot")
	}
}

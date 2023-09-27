package payload

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"

	"github.com/sirupsen/logrus"
)

type Meta struct {
	Tag  RecordTag
	Type RecordType
}

type Dependency uint8

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
	RecordTagName RecordTag = 1
	// Architecture of the package
	RecordTagArchitecture RecordTag = 2
	// Version of the package
	RecordTagVersion RecordTag = 3
	// Summary of the package
	RecordTagSummary RecordTag = 4
	// Description of the package
	RecordTagDescription RecordTag = 5
	// Homepage for the package
	RecordTagHomepage RecordTag = 6
	// ID for the source package, used for grouping
	RecordTagSourceID RecordTag = 7
	// Runtime dependencies
	RecordTagDepends RecordTag = 8
	// Provides some capability or name
	RecordTagProvides RecordTag = 9
	// Conflicts with some capability or name
	RecordTagConflicts RecordTag = 10
	// Release number for the package
	RecordTagRelease RecordTag = 11
	// SPDX license identifier
	RecordTagLicense RecordTag = 12
	// Currently recorded build number
	RecordTagBuildRelease RecordTag = 13
	// Repository index specific (relative URI)
	RecordTagPackageURI RecordTag = 14
	// Repository index specific (Package hash)
	RecordTagPackageHash RecordTag = 15
	// Repository index specific (size on disk)
	RecordTagPackageSize RecordTag = 16
	// A Build Dependency
	RecordTagBuildDepends RecordTag = 17
	// Upstream URI for the source
	RecordTagSourceURI RecordTag = 18
	// Relative path for the source within the upstream URI
	RecordTagSourcePath RecordTag = 19
	// Ref/commit of the upstream source
	RecordTagSourceRef RecordTag = 20
)

type RecordType uint8

const (
	RecordTypeUnknown    RecordType = 0
	RecordTypeInt8       RecordType = 1
	RecordTypeUint8      RecordType = 2
	RecordTypeInt16      RecordType = 3
	RecordTypeUint16     RecordType = 4
	RecordTypeInt32      RecordType = 5
	RecordTypeUint32     RecordType = 6
	RecordTypeInt64      RecordType = 7
	RecordTypeUint64     RecordType = 8
	RecordTypeString     RecordType = 9
	RecordTypeDependency RecordType = 10
	RecordTypeProvider   RecordType = 11
)

type MetaRecord struct {
	Length     uint32
	RecordTag  RecordTag
	RecordType RecordType
	Padding    byte
}

func DecodeMetaPayload(payload []byte, records int) error {
	rawReader := bytes.NewReader(payload)
	reader := bufio.NewReader(rawReader)
	offset := 0
	for i := 0; i < records; i++ {
		record := MetaRecord{}

		err := binary.Read(bytes.NewReader(payload[offset:offset+8]), binary.BigEndian, &record)
		if err != nil {
			return err
		}
		offset = offset + 8
		logrus.Printf("Payload %d, Length %d, Tag %s, Type %s", i, record.Length, record.RecordTag.String(), record.RecordType.String())
		_, err = rawReader.Seek(8, io.SeekCurrent)
		if err != nil {
			return err
		}
		switch record.RecordType {
		case RecordTypeString:
			output, err := reader.ReadString('\x00')
			if err != nil {
				return err
			}
			logrus.Printf("Output: %s", output)
		default:
			_, err = rawReader.Seek(int64(record.Length), io.SeekCurrent)
			if err != nil {
				return err
			}
		}
		offset = offset + int(record.Length)
	}
	return nil
}

package payload

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

//go:generate -command stringer go run golang.org/x/tools/cmd/stringer
//go:generate stringer -type RecordType,RecordTag,Kind,Compression,Dependency

type Dependency uint8

const (
	DependencyPackageName Dependency = iota
	DependencySharedLibrary
	DependencyPkgConfig
	DependencyInterpreter
	DependencyCMake
	DependencyPython
	DependencyBinary
	DependencySystemBinary
	DependencyPkgConfig32
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

func ReadIntegerData[T any](input io.Reader) (T, error) {
	var output T
	err := binary.Read(input, binary.BigEndian, &output)
	if err != nil {
		return output, err
	}
	return output, nil
}

func ReadDependsProvides(buf *bytes.Buffer) (string, error) {
	depType, err := ReadIntegerData[uint8](buf)
	if err != nil {
		return "", err
	}
	depends, err := buf.ReadString('\x00')
	if err != nil {
		return "", err
	}
	return wrapDependency(Dependency(depType), depends), nil
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

		fmt.Printf("%-15s : %v\n", strings.TrimLeft(record.RecordTag.String(), "RecordTag"), data)
	}
	return nil
}

func wrapDependency(depType Dependency, name string) string {
	switch depType {
	case DependencyPackageName:
		return name
	case DependencySharedLibrary:
		return name
	case DependencyPkgConfig:
		return fmt.Sprintf("pkgconfig(%s)", name)
	case DependencyInterpreter:
		return fmt.Sprintf("interpreter(%s)", name)
	case DependencyCMake:
		return fmt.Sprintf("cmake(%s)", name)
	case DependencyPython:
		return fmt.Sprintf("python(%s)", name)
	case DependencyBinary:
		return fmt.Sprintf("binary(%s)", name)
	case DependencySystemBinary:
		return fmt.Sprintf("system_binary(%s)", name)
	case DependencyPkgConfig32:
		return fmt.Sprintf("pkgconfig32(%s)", name)
	}
	return name
}

func switchstuff(buf *bytes.Buffer, recordType RecordType) (any, error) {
	switch recordType {
	case RecordTypeInt8:
		return ReadIntegerData[int8](buf)
	case RecordTypeUint8:
		return ReadIntegerData[uint8](buf)
	case RecordTypeInt16:
		return ReadIntegerData[int16](buf)
	case RecordTypeUint16:
		return ReadIntegerData[uint16](buf)
	case RecordTypeInt32:
		return ReadIntegerData[int32](buf)
	case RecordTypeUint32:
		return ReadIntegerData[uint32](buf)
	case RecordTypeInt64:
		return ReadIntegerData[int64](buf)
	case RecordTypeUint64:
		return ReadIntegerData[uint64](buf)
	case RecordTypeString:
		return buf.ReadString('\x00')
	case RecordTypeDependency:
		return ReadDependsProvides(buf)
	case RecordTypeProvider:
		return ReadDependsProvides(buf)
	default:
		return nil, errors.Errorf("Unknown RecordType: %s", recordType.String())
	}
}

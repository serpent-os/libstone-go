// SPDX-FileCopyrightText: 2024 Serpent OS Developers
// SPDX-License-Identifier: MPL-2.0

package stone1

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"

	"github.com/serpent-os/libstone/internal/readers"
	"github.com/zeebo/xxh3"
)

// Record is a payload record. There are multiple
// kinds of record: a type conversion is required to inspect the content.
type Record interface {
	// Kind returns the kind of this record.
	Kind() RecordKind
	decode(src io.Reader) error
}

type AttributeRecord struct {
	Key   []byte
	Value []byte
}

// Kind returns the kind of this record.
func (r AttributeRecord) Kind() RecordKind {
	return Attributes
}

func (r *AttributeRecord) decode(src io.Reader) error {
	var lengths [8 + 8]byte
	_, err := io.ReadFull(src, lengths[:])
	if err != nil {
		return err
	}
	wlk := readers.ByteWalker(lengths[:])
	keyLen := wlk.Uint64()
	valLen := wlk.Uint64()

	r.Key = make([]byte, keyLen)
	_, err = io.ReadFull(src, r.Key)
	if err != nil {
		return err
	}
	r.Value = make([]byte, valLen)
	_, err = io.ReadFull(src, r.Value)
	if err != nil {
		return err
	}
	return nil
}

// IndexRecord records offsets to unique files within the content when decompressed.
// This is used to split the file into the content store on disk before promoting
// to a transaction.
type IndexRecord struct {
	// Start is the index where the content starts.
	Start uint64
	// End is the index where the content ends.
	End uint64
	// Hash is the XXH3_128 hash of the content.
	Hash xxh3.Uint128
}

// Kind returns the kind of this record.
func (r IndexRecord) Kind() RecordKind {
	return Index
}

func (r *IndexRecord) decode(src io.Reader) error {
	var data [8 + 8 + 16]byte
	_, err := io.ReadFull(src, data[:])
	if err != nil {
		return err
	}
	wlk := readers.ByteWalker(data[:])
	*r = IndexRecord{
		Start: wlk.Uint64(),
		End:   wlk.Uint64(),
		Hash: xxh3.Uint128{
			Hi: wlk.Uint64(),
			Lo: wlk.Uint64(),
		},
	}
	return nil
}

type MetaTag uint16

const (
	// Name of the package
	Name MetaTag = iota + 1
	// Architecture is the architecture of the package.
	Architecture
	// Version is the version of the package.
	Version
	// Summary is the succint description of the package.
	Summary
	// Description is the description of the package.
	Description
	// Homepage is the homepage URL of the package.
	Homepage
	// SourceID is the ID of the source package used for grouping.
	SourceID // Source ID
	// Depends is one dependency of the package.
	Depends
	// Provides is one capability, or the name, of the package.
	Provides
	// Conflicts is one capability, or name, conflicting with this package.
	Conflicts
	// Release is the release number of the package.
	Release
	// SPDX lists the SPDX license identifiers of the package.
	License
	// BuildRelease is the currently recorded build number of the package.
	BuildRelease // Build release
	// PackageURI is the URI of the package.
	PackageURI // Package URI
	// PackageHash is the hash sum of the package.
	PackageHash // Package hash
	// PackageSize is the size of the package.
	PackageSize // Package size
	// Depends is one build-time dependency of the package.
	BuildDepends // Build dependencies
	// SourceURI is the URI of the source of the package.
	SourceURI // Source URI
	// SourcePath is the relative path for the source within the upstream URI.
	SourcePath // Source path
	// SourceRef is the ref (or commit) of the upstream source.
	SourceRef // Source ref
)

type MetaFieldKind uint8

const (
	Int8MetaField MetaFieldKind = iota + 1
	Uint8MetaField
	Int16MetaField
	Uint16MetaField
	Int32MetaField
	Uint32MetaField
	Int64MetaField
	Uint64MetaField
	StringMetaField
	DependencyMetaField
	ProviderMetaField
)

type MetaField struct {
	Kind  MetaFieldKind
	Value any
}

func (m MetaField) String() string {
	return fmt.Sprint(m.Value)
}

func (mv MetaField) size() int {
	switch mv.Kind {
	case Int8MetaField, Uint8MetaField:
		return 1
	case Int16MetaField, Uint16MetaField:
		return 2
	case Int32MetaField, Uint32MetaField:
		return 4
	case Int64MetaField, Uint64MetaField:
		return 8
	case StringMetaField:
		return len(mv.Value.(string))
	case DependencyMetaField, ProviderMetaField:
		// TODO: should cast to a DependencyValue.
		return len(mv.Value.(string))
	default:
		panic("unknown MetaKind value")
	}
}

type DependencyKind uint8

const (
	// PackageName is the name of the package.
	PackageName DependencyKind = iota // name
	// SharedLibary is a soname based dependency.
	SharedLibary // soname
	// PkgConfig is a pkgconfig `.pc` based dependency.
	PkgConfig // pkgconfig
	// Interpreter is a special interpreter (PT_INTERP/etc) to run the binaries.
	Interpreter // interpreter
	// CMake is a CMake module.
	CMake // cmake
	// Python is a Python module.
	Python // python
	// Binary is a binary in /usr/bin.
	BinaryDep // binary
	// SystemBinary is a binary in /usr/sbin.
	SystemBinary // sysbinary
	// PkgConfig32 is a emul32-compatible pkgconfig .pc dependency (contained in lib32/*.pc).
	PkgConfig32 // pkgconfig32
)

type Dependency struct {
	Kind DependencyKind
	Name string
}

func (d Dependency) String() string {
	return fmt.Sprintf("%s(%s)", d.Kind, d.Name)
}

type MetaRecord struct {
	Tag   MetaTag
	Field MetaField
}

// Kind returns the kind of this record.
func (r MetaRecord) Kind() RecordKind {
	return Meta
}

func (r *MetaRecord) decode(src io.Reader) error {
	var header [4 + 2 + 1 + 1]byte
	_, err := io.ReadFull(src, header[:])
	if err != nil {
		return err
	}

	wlk := readers.ByteWalker(header[:])
	length := wlk.Uint32()
	r.Tag = MetaTag(wlk.Uint16())
	r.Field.Kind = MetaFieldKind(wlk.Uint8())
	wlk.Uint8() // Skip padding.

	buf := make([]byte, length)
	_, err = io.ReadFull(src, buf)
	if err != nil {
		return err
	}
	switch r.Field.Kind {
	case Int8MetaField:
		r.Field.Value = int8(buf[0])
	case Uint8MetaField:
		r.Field.Value = buf[0]
	case Int16MetaField:
		r.Field.Value = int16(readers.ByteOrder.Uint16(buf))
	case Uint16MetaField:
		r.Field.Value = readers.ByteOrder.Uint16(buf)
	case Int32MetaField:
		r.Field.Value = int32(readers.ByteOrder.Uint32(buf))
	case Uint32MetaField:
		r.Field.Value = readers.ByteOrder.Uint32(buf)
	case Int64MetaField:
		r.Field.Value = int64(readers.ByteOrder.Uint64(buf))
	case Uint64MetaField:
		r.Field.Value = readers.ByteOrder.Uint64(buf)
	case StringMetaField:
		r.Field.Value = trimTerminator(buf)
	case DependencyMetaField, ProviderMetaField:
		r.Field.Value = Dependency{
			Kind: DependencyKind(buf[0]),
			Name: trimTerminator(buf[1:]),
		}
	}
	return nil
}

type FileType uint8

const (
	// Regular is a regular file.
	Regular FileType = iota + 1
	// Symlink is a symbolic link (source, target pair).
	Symlink
	// Directory is a directory node.
	Directory
	// CharacterDevice is a character device.
	CharacterDevice // Character device
	// BlockDevice is a block device.
	BlockDevice // Block device
	// FIFO is a FIFO node.
	FIFO
	// Socket is a UNIX socket.
	Socket
)

type Entry struct {
	FileType FileType
	value    any
}

func (e Entry) Source() []byte {
	switch e.FileType {
	case Regular:
		hashAndTarget := e.value.(tuple[xxh3.Uint128, string])
		hash := hashAndTarget.val1.Bytes()
		return hash[:]
	case Symlink:
		sourceAndTarget := e.value.(tuple[string, string])
		return []byte(sourceAndTarget.val1)
	case Directory,
		CharacterDevice,
		BlockDevice,
		FIFO,
		Socket:
		return nil
	default:
		panic("unknown value of FileType")
	}
}

func (e Entry) Target() []byte {
	switch e.FileType {
	case Regular:
		hashAndTarget := e.value.(tuple[xxh3.Uint128, string])
		return []byte(hashAndTarget.val2)
	case Symlink:
		sourceAndTarget := e.value.(tuple[string, string])
		return []byte(sourceAndTarget.val2)
	case Directory,
		CharacterDevice,
		BlockDevice,
		FIFO,
		Socket:
		target := e.value.(string)
		return []byte(target)
	default:
		panic("unknown value of FileType")
	}
}

// LayoutRecord contains information about a
// file that should be written to the mass memory.
type LayoutRecord struct {
	// UID is the UNIX UID.
	UID uint32
	// GID is the UNIX GID.
	GID uint32
	// Mode is file's mode.
	Mode fs.FileMode
	Tag  uint32
	// Entry is the kind of file, with source
	// and target paths where necessary.
	Entry Entry
}

// Kind returns the kind of this record.
func (r LayoutRecord) Kind() RecordKind {
	return Layout
}

func (r *LayoutRecord) decode(src io.Reader) error {
	var header [4 + 4 + 4 + 4 + 2 + 2 + 1 + 11]byte
	_, err := io.ReadFull(src, header[:])
	if err != nil {
		return err
	}

	wlk := readers.ByteWalker(header[:])
	*r = LayoutRecord{
		UID:  wlk.Uint32(),
		GID:  wlk.Uint32(),
		Mode: fs.FileMode(wlk.Uint32()),
		Tag:  wlk.Uint32(),
	}
	srcLen := wlk.Uint16()
	tgtLen := wlk.Uint16()
	r.Entry.FileType = FileType(wlk.Uint8())
	wlk.Ahead(11) // Skip padding.

	buf := make([]byte, srcLen+tgtLen)
	_, err = io.ReadFull(src, buf)
	if err != nil {
		return err
	}
	switch r.Entry.FileType {
	case Regular:
		r.Entry.value = tuple[xxh3.Uint128, string]{
			val1: xxh3.Uint128{
				Hi: readers.ByteOrder.Uint64(buf[:srcLen/2]),
				Lo: readers.ByteOrder.Uint64(buf[srcLen/2 : srcLen]),
			},
			val2: trimTerminator(buf[srcLen:]),
		}
	case Symlink:
		r.Entry.value = tuple[string, string]{
			val1: trimTerminator(buf[:srcLen]),
			val2: trimTerminator(buf[srcLen:]),
		}
	case Directory,
		CharacterDevice,
		BlockDevice,
		FIFO,
		Socket:
		r.Entry.value = trimTerminator(buf[:tgtLen])
	default:
		panic("unknown value of FileType")
	}
	return nil
}

type ContentRecord struct {
	Data *io.LimitedReader
}

// Kind returns the kind of this record.
func (r ContentRecord) Kind() RecordKind {
	return Content
}

func (r *ContentRecord) decode(src io.Reader) error {
	reader, ok := src.(*io.LimitedReader)
	if !ok {
		panic("src is not a *io.LimitedReader")
	}
	r.Data = reader
	return nil
}

// tuple mimics the tuple type from other languages.
type tuple[T1, T2 any] struct {
	val1 T1
	val2 T2
}

func trimTerminator(str []byte) string {
	return string(bytes.TrimSuffix(str, []byte{0}))
}

//go:generate go run golang.org/x/tools/cmd/stringer@v0.18.0 -linecomment -type=MetaTag,DependencyKind,FileType -output record_enumstring.go

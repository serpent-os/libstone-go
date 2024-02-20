package stone1

import (
	"github.com/serpent-os/libstone/internal/readers"
)

// RecordKind is the kind of the payload's records.
type RecordKind uint8

const (
	// Meta indicates a [MetaRecord].
	Meta RecordKind = iota + 1
	// Content indicates a file.
	Content
	// Layout indicates a [LayoutRecord].
	Layout
	// Index indicates an [IndexRecord].
	Index
	// Attributes indicates an attribute store.
	Attributes
)

// Compression is the compression method of the archive content.
type Compression uint8

const (
	// Uncompressed indicates an uncompressed content.
	Uncompressed Compression = iota + 1
	// ZSTD indicates a compressed content using zstd.
	ZSTD
)

// Header is the payload's header.
type Header struct {
	// StoredSize is the compressed size, in bytes, of the payload.
	// If the payload is not compressed, StoredSize is equal to PlainSize.
	StoredSize uint64
	// PlainSize is the uncompressed size, in bytes, of the payload.
	PlainSize uint64
	// Checksum is the payload's checksum.
	Checksum uint64
	// NumRecords is the number of records contained in the payload.
	NumRecords uint32
	// Version is the version of the payload data format.
	Version uint16
	// Kind is the kind of payload's records.
	Kind RecordKind
	// Compression is the compression used for the payload.
	Compression Compression
}

const (
	headerLen = 32
)

func newHeader(data [headerLen]byte) Header {
	wlk := readers.ByteWalker(data[:])
	return Header{
		StoredSize:  wlk.Uint64(),
		PlainSize:   wlk.Uint64(),
		Checksum:    wlk.Uint64(),
		NumRecords:  wlk.Uint32(),
		Version:     wlk.Uint16(),
		Kind:        RecordKind(wlk.Uint8()),
		Compression: Compression(wlk.Uint8()),
	}
}

//go:generate go run golang.org/x/tools/cmd/stringer@v0.18.0 -linecomment -type=RecordKind -output payload_enumstring.go

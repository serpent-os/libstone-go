package payload

type Kind uint8

const (
	KindMeta Kind = iota + 1
	KindContent
	KindLayout
	KindIndex
	KindAttributes
	KindDumb
)

type Meta struct {
	Tag  Tag
	Kind Kind
}

type Compression uint8

const (
	None Compression = 1
	Zstd             = 2
)

type PayloadHeader struct {
	StoredSize  uint64
	PlainSize   uint64
	Checksum    [8]uint8
	NumRecords  uint32
	Version     uint16
	Kind        Kind
	Compression Compression
}

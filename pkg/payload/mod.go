package payload

import (
	"bytes"
	"encoding/binary"

	"github.com/sirupsen/logrus"
)

type Kind uint8

const (
	KindMeta Kind = iota + 1
	KindContent
	KindLayout
	KindIndex
	KindAttributes
	KindDumb
)

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

func ReadPayloadHeader(headerData [32]byte) (*PayloadHeader, error) {
	payloadHeaderHeader := PayloadHeader{}
	r := bytes.NewReader(headerData[:])
	err := binary.Read(r, binary.BigEndian, &payloadHeaderHeader)
	if err != nil {
		return nil, err
	}

	return &payloadHeaderHeader, nil
}

func (p PayloadHeader) Print() {
	logrus.Printf("Payload kind: %d", p.Kind)
	logrus.Printf("Payload Compression: %d", p.Compression)
	logrus.Printf("Payload version: %d", p.Version)
	logrus.Printf("Payload records: %d", p.NumRecords)
	logrus.Printf("Payload stored size: %d", p.StoredSize)
	logrus.Printf("Payload plain size: %d", p.PlainSize)
}

const SIZE_PAYLOAD_HEADER = 8 + 8 + 8 + 4 + 2 + 1 + 1

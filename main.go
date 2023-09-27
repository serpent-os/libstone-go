package main

import (
	"fmt"
	"log"
	"os"

	"github.com/der-eismann/libstone/pkg/header"
	"github.com/sirupsen/logrus"
)

type DataHeader struct {
	NumPayloads    uint16
	IntegrityCheck [21]byte
	FileType       byte
}

const (
	PayloadFileTypeRegular = iota + 1
	PayloadFileTypeSymlink
	PayloadFileTypeDirectory
	PayloadFileTypeCharacterDevice
	PayloadFileTypeBlockDevice
	PayloadFileTypeFifo
	PayloadFileTypeDocket
)

type LayoutStruct struct {
	UID  uint32
	GID  uint32
	Mode uint32
	Tag  uint32
}

type Kind uint8

const (
	Meta       Kind = 1
	Content    Kind = 2
	Layout     Kind = 3
	Index      Kind = 4
	Attributes Kind = 5
	Dumb       Kind = 6
)

type Compression uint8

const (
	None Compression = 1
	Zstd Compression = 2
)

type PayloadDings struct {
	RealSize         uint64
	DecompressedSize uint64
	Checksum         [8]uint8
	PayloadRecords   uint32
	PayloadVersion   uint16
	PayloadType      Kind
	CompressionType  Compression
}

func main() {
	file, err := os.Open("bash-completion-2.11-1-1-x86_64.stone") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	data := [32]byte{}
	count, err := file.Read(data[:])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("read %d bytes: %q\n", count, data[:count])

	// header := AgnosticHeader{
	// 	Magic:   data[0:4],
	// 	Data:    data[4:28],
	// 	Version: data[28:32],
	// }

	// fmt.Printf("Read struct: %#v\n", header)

	// var numPayloads uint16
	// r := bytes.NewReader(data[4:6])
	// err = binary.Read(r, binary.BigEndian, &numPayloads)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// header2 := DataHeader{
	// 	NumPayloads:    numPayloads,
	// 	IntegrityCheck: data[6:23],
	// 	FileType:       data[23],
	// }

	// fmt.Printf("Stone contains %d files\n", numPayloads)

	// var dings PayloadDings
	// r = bytes.NewReader(data[32:64])
	// err = binary.Read(r, binary.BigEndian, &dings)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("Size: %d, Size2: %d, Records: %d, Type: %d, Compression: %d\n", dings.RealSize, dings.DecompressedSize, dings.PayloadRecords, dings.PayloadType, dings.CompressionType)

	header, err := header.ReadHeader(data)
	logrus.Printf("Header decoded:")
	logrus.Printf("- Number of Payloads: %d", header.Data.NumPayloads)
}

const SIZE_PAYLOAD_HEADER = 8 + 8 + 8 + 4 + 2 + 1 + 1

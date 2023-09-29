package main

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/der-eismann/libstone/pkg/header"
	"github.com/der-eismann/libstone/pkg/payload"
	"github.com/der-eismann/libstone/pkg/zstd"
	"github.com/sirupsen/logrus"
)

func main() {
	file, err := os.Open("bash-completion-2.11-1-1-x86_64.stone") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	data := [32]byte{}
	_, err = file.Read(data[:])
	if err != nil {
		log.Fatal(err)
	}

	header, err := header.ReadHeader(data)
	if err != nil {
		log.Fatal(err)
	}
	logrus.Printf("Header decoded:")
	logrus.Printf("- Number of Payloads: %d", header.Data.NumPayloads)
	logrus.Printf("- FileType: %d", header.Data.FileType)

	for i := 0; i < int(header.Data.NumPayloads); i++ {
		_, err = file.Read(data[:])
		if err != nil {
			log.Fatal(err)
		}
		payloadheader, err := payload.ReadPayloadHeader(data)
		if err != nil {
			log.Fatal(err)
		}
		payloadheader.Print()

		pos, err := file.Seek(0, io.SeekCurrent)
		if err != nil {
			log.Fatal(err)
		}

		sectionReader := io.NewSectionReader(file, pos, int64(payloadheader.StoredSize))

		decompdata := make([]byte, 0, payloadheader.PlainSize)
		writer := bytes.NewBuffer(decompdata)

		decomp, err := zstd.Decompress(sectionReader, writer)
		if err != nil {
			log.Fatal(err)
		}
		logrus.Printf("Bytes copied: %d", decomp)
		logrus.Printf("%#v", writer.Bytes())
		err = payload.DecodeMetaPayload(writer.Bytes(), int(payloadheader.NumRecords))
		if err != nil {
			log.Fatal(err)
		}

		_, err = file.Seek(int64(payloadheader.StoredSize), io.SeekCurrent)
		if err != nil {
			log.Fatal(err)
		}

		//file.Seek(int64(payloadheader.StoredSize), io.SeekCurrent)
		break
	}
}

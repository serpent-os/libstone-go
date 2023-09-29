package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/der-eismann/libstone/pkg/header"
	"github.com/der-eismann/libstone/pkg/payload"
	"github.com/der-eismann/libstone/pkg/zstd"
)

const FILE_NAME = "bash-completion-2.11-1-1-x86_64.stone"

func main() {
	fmt.Printf("Archive: %s\n", FILE_NAME)
	file, err := os.Open(FILE_NAME) // For read access.
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

		_, err = zstd.Decompress(sectionReader, writer)
		if err != nil {
			log.Fatal(err)
		}
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

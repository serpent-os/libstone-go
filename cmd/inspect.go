package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/der-eismann/libstone/pkg/header"
	"github.com/der-eismann/libstone/pkg/payload"
	"github.com/der-eismann/libstone/pkg/zstd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Inspect(ctx context.Context, cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		logrus.Fatal("One stone file as argument required")
	}

	fmt.Printf("Archive: %s\n", args[0])
	file, err := os.Open(args[0]) // For read access.
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

		payloadData := []byte{}

		if payloadheader.Compression == payload.Zstd {
			decompdata := make([]byte, 0, payloadheader.PlainSize)
			writer := bytes.NewBuffer(decompdata)
			_, err = zstd.Decompress(sectionReader, writer)
			if err != nil {
				log.Fatal(err)
			}
			payloadData = writer.Bytes()
		} else {
			_, err = sectionReader.Read(payloadData)
			if err != nil {
				log.Fatal(err)
			}
		}

		switch payloadheader.Kind {
		case payload.KindMeta:
			err = payload.DecodeMetaPayload(payloadData, int(payloadheader.NumRecords))
		case payload.KindLayout:
			err = payload.DecodeLayoutPayload(payloadData, int(payloadheader.NumRecords))
		case payload.KindIndex:
			err = payload.DecodeIndexPayload(payloadData, int(payloadheader.NumRecords))
		default:
			continue
		}
		if err != nil {
			log.Fatal(err)
		}

		_, err = file.Seek(int64(payloadheader.StoredSize), io.SeekCurrent)
		if err != nil {
			log.Fatal(err)
		}
	}
}

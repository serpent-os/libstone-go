package cmd

import (
	"bytes"
	"context"
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

	var pos int64

	file, err := os.Open(args[0])
	if err != nil {
		logrus.Fatalf("Failed to open file: %s", err)
	}

	packageHeader, err := header.ReadHeader(io.NewSectionReader(file, 0, 32))
	if err != nil {
		logrus.Fatalf("Failed to read package header: %s", err)
	}

	pos += 32

	for i := 0; i < int(packageHeader.Data.NumPayloads); i++ {
		payloadheader, err := payload.ReadPayloadHeader(io.NewSectionReader(file, pos, 32))
		if err != nil {
			logrus.Fatalf("Failed to read payload header: %s", err)
		}
		payloadheader.Print()

		pos += 32

		sectionReader := io.NewSectionReader(file, pos, int64(payloadheader.StoredSize))

		pos += int64(payloadheader.StoredSize)

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
	}
}

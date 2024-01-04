package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/der-eismann/libstone/pkg/header"
	"github.com/der-eismann/libstone/pkg/payload"
	"github.com/klauspost/compress/zstd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Inspect(ctx context.Context, cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		logrus.Fatal("One stone file as argument required")
	}

	var pos int64

	absPath, err := filepath.Abs(args[0])
	if err != nil {
		logrus.Fatalf("Failed to get absolute path: %s", err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		logrus.Fatalf("Failed to open file: %s", err)
	}

	fmt.Printf("\"%s\" = stone container version V1\n", absPath)

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
		//payloadheader.Print()

		pos += 32

		payloadReader, err := getCompressionReader(file, payloadheader.Compression, pos, int64(payloadheader.StoredSize))
		if err != nil {
			logrus.Fatalf("Failed to get compression reader: %s", err)
		}

		pos += int64(payloadheader.StoredSize)

		switch payloadheader.Kind {
		case payload.KindMeta:
			err = payload.PrintMetaPayload(payloadReader, int(payloadheader.NumRecords))
		case payload.KindLayout:
			err = payload.PrintLayoutPayload(payloadReader, int(payloadheader.NumRecords))
		// case payload.KindIndex:
		// 	err = payload.PrintIndexPayload(payloadReader, int(payloadheader.NumRecords))
		default:
			continue
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getCompressionReader(r io.ReaderAt, compressionType payload.Compression, offset, length int64) (io.Reader, error) {
	switch compressionType {
	case payload.CompressionNone:
		return io.NewSectionReader(r, offset, length), nil
	case payload.CompressionZstd:
		return zstd.NewReader(io.NewSectionReader(r, offset, length))
	}
	return nil, errors.New("Unknown compression type")
}

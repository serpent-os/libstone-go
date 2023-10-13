package payload

import (
	"encoding/binary"
	"fmt"
	"io"
)

type IndexEntry struct {
	Start  uint64
	End    uint64
	Digest [16]uint8
}

func PrintIndexPayload(r io.Reader, records int) error {
	for i := 0; i < records; i++ {
		record := IndexEntry{}

		err := binary.Read(r, binary.BigEndian, &record)
		if err != nil {
			return err
		}

		fmt.Printf("  - %x [size: %s]\n", record.Digest, formatBytes(record.End-record.Start))
	}
	return nil
}

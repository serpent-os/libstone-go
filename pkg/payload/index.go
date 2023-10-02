package payload

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type IndexEntry struct {
	Start  uint64
	End    uint64
	Digest [16]uint8
}

func DecodeIndexPayload(payload []byte, records int) error {
	reader := bytes.NewBuffer(payload)
	for i := 0; i < records; i++ {
		record := IndexEntry{}

		err := binary.Read(reader, binary.BigEndian, &record)
		if err != nil {
			return err
		}

		fmt.Printf("  - %x [size: %s]\n", record.Digest, formatBytes(record.End-record.Start))
	}
	return nil
}

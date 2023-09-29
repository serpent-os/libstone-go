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

		fmt.Printf("  - %x [size: %s]\n", record.Digest, ByteFormatIEC(record.End-record.Start))
		//- 3c005061e2d565b469e9abdfe6478cfe [size:    74.54 KiB]
	}
	return nil
}

func ByteFormatIEC(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%9d   B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%9.2f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

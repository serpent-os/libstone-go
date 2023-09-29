package payload

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type FileType uint8

const (
	FileTypeRegular FileType = iota + 1
	FileTypeSymlink
	FileTypeDirectory
	FileTypeCharacterDevice
	FileTypeBlockDevice
	FileTypeFifo
	FileTypeSocket
)

type LayoutEntry struct {
	UID          uint32
	GID          uint32
	Mode         uint32
	Tag          uint32
	SourceLength uint16
	TargetLength uint16
	FileType     FileType
	Padding      [11]byte
}

func DecodeLayoutPayload(payload []byte, records int) error {
	reader := bytes.NewBuffer(payload)
	for i := 0; i < records; i++ {
		record := LayoutEntry{}

		err := binary.Read(reader, binary.BigEndian, &record)
		if err != nil {
			return err
		}

		switch record.FileType {
		case FileTypeRegular:
			pt1, err := ReadIntegerData[uint64](reader)
			if err != nil {
				return err
			}
			pt2, err := ReadIntegerData[uint64](reader)
			if err != nil {
				return err
			}
			source, err := reader.ReadString('\x00')
			if err != nil {
				return err
			}
			fmt.Printf("  - /usr/%s -> %x%x [%s]\n", source, pt1, pt2, strings.TrimLeft(record.FileType.String(), "FileType"))
		case FileTypeSymlink:
			target, err := reader.ReadString('\x00')
			if err != nil {
				return err
			}
			source, err := reader.ReadString('\x00')
			if err != nil {
				return err
			}
			fmt.Printf("  - /usr/%s -> %s [%s]\n", source, target, strings.TrimLeft(record.FileType.String(), "FileType"))
		default:
			source, err := reader.ReadString('\x00')
			if err != nil {
				return err
			}
			fmt.Printf("  - /usr/%s [%s]\n", source, strings.TrimLeft(record.FileType.String(), "FileType"))
		}

		//fmt.Printf("  - %s -> %s [%s]", source, target, strings.TrimLeft(record.FileType.String(), "FileType"))
		//  - /usr/share/bash-completion/bash_completion -> 3c005061e2d565b469e9abdfe6478cfe [Regular]

	}
	return nil
}

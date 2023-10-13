package payload

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
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

func PrintLayoutPayload(r io.Reader, records int) error {
	bufferedReader := bufio.NewReader(r)
	for i := 0; i < records; i++ {
		record := LayoutEntry{}

		err := binary.Read(bufferedReader, binary.BigEndian, &record)
		if err != nil {
			return err
		}

		switch record.FileType {
		case FileTypeRegular:
			pt1, err := ReadIntegerData[uint64](bufferedReader)
			if err != nil {
				return err
			}
			pt2, err := ReadIntegerData[uint64](bufferedReader)
			if err != nil {
				return err
			}
			source, err := bufferedReader.ReadBytes('\x00')
			if err != nil {
				return err
			}
			fmt.Printf("  - /usr/%s -> %x%x [%s]\n", source[:len(source)-1], pt1, pt2, strings.TrimLeft(record.FileType.String(), "FileType"))
		case FileTypeSymlink:
			target, err := bufferedReader.ReadBytes('\x00')
			if err != nil {
				return err
			}
			source, err := bufferedReader.ReadBytes('\x00')
			if err != nil {
				return err
			}
			fmt.Printf("  - /usr/%s -> %s [%s]\n", source[:len(source)-1], target[:len(target)-1], strings.TrimLeft(record.FileType.String(), "FileType"))
		default:
			source, err := bufferedReader.ReadString('\x00')
			if err != nil {
				return err
			}
			fmt.Printf("  - /usr/%s [%s]\n", source[:len(source)-1], strings.TrimLeft(record.FileType.String(), "FileType"))
		}
	}
	return nil
}

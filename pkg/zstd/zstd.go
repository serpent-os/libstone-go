package zstd

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

// func Decompress(in []byte, decompressedSize uint64) ([]byte, error) {
// 	bytesReader := bytes.NewReader(in)
// 	d, err := zstd.NewReader(bytesReader)
// 	if err != nil {
// 		return []byte{}, err
// 	}
// 	defer d.Close()

//		buf := new(bytes.Buffer)
//		n, err := io.Copy(buf, d)
//		if err != nil {
//			return []byte{}, err
//		}
//		if uint64(n) != decompressedSize {
//			return []byte{}, errors.New("written size does not matched decompressed size")
//		}
//		return buf.Bytes(), nil
//	}

func Decompress(in io.Reader, out io.Writer) (int64, error) {
	d, err := zstd.NewReader(in)
	if err != nil {
		return 0, err
	}
	defer d.Close()

	// Copy content...
	n, err := io.Copy(out, d)
	return n, err
}

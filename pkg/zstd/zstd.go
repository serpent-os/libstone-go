package zstd

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

func Decompress(in io.Reader, out io.Writer) (int64, error) {
	d, err := zstd.NewReader(in)
	if err != nil {
		return 0, err
	}
	defer d.Close()

	n, err := io.Copy(out, d)
	return n, err
}

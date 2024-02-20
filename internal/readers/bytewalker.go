package readers

import (
	"encoding/binary"
)

var (
	ByteOrder = binary.BigEndian
)

type ByteWalker []byte

func (r *ByteWalker) Ahead(n int) []byte {
	val := (*r)[:n]
	*r = (*r)[n:]
	return val
}

func (r *ByteWalker) Uint8() uint8 {
	val := (*r)[0]
	*r = (*r)[1:]
	return val
}

func (r *ByteWalker) Uint16() uint16 {
	val := ByteOrder.Uint16(*r)
	*r = (*r)[2:]
	return val
}

func (r *ByteWalker) Uint32() uint32 {
	val := ByteOrder.Uint32(*r)
	*r = (*r)[4:]
	return val
}

func (r *ByteWalker) Uint64() uint64 {
	val := ByteOrder.Uint64(*r)
	*r = (*r)[8:]
	return val
}

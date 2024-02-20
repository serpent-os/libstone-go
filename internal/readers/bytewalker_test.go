package readers_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/serpent-os/libstone/internal/readers"
)

const (
	aheadDistance = 5
)

var (
	testData []byte
)

func init() {
	testData = make([]byte, 8+4+2+1+aheadDistance)
	for i := range testData {
		testData[i] = byte(i + 1)
	}
}

func TestAhead(t *testing.T) {
	wlk := readers.ByteWalker(testData)
	testAhead(t, &wlk, testData)
}

func TestUint8(t *testing.T) {
	wlk := readers.ByteWalker(testData)
	testUint8(t, &wlk, testData)
}

func TestUint16(t *testing.T) {
	wlk := readers.ByteWalker(testData)
	testUint16(t, &wlk, testData)
}

func TestUint32(t *testing.T) {
	wlk := readers.ByteWalker(testData)
	testUint32(t, &wlk, testData)
}

func TestUint64(t *testing.T) {
	wlk := readers.ByteWalker(testData)
	testUint64(t, &wlk, testData)
}

func TestIsWalking(t *testing.T) {
	wlk := readers.ByteWalker(testData)
	data := testData
	testAhead(t, &wlk, data)
	testUint8(t, &wlk, data[aheadDistance:])
	testUint16(t, &wlk, data[aheadDistance+1:])
	testUint32(t, &wlk, data[aheadDistance+1+2:])
	testUint64(t, &wlk, data[aheadDistance+1+2+4:])
}

func testAhead(t *testing.T, wlk *readers.ByteWalker, data []byte) {
	expect := data[:aheadDistance]
	obtain := wlk.Ahead(aheadDistance)
	if !bytes.Equal(obtain, expect) {
		t.Fatalf("expected ahead slice %v. Got %v", expect, obtain)
	}
}

func testUint8(t *testing.T, wlk *readers.ByteWalker, data []byte) {
	expect := data[0]
	obtain := wlk.Uint8()
	if obtain != expect {
		t.Fatalf("expected uint8 %d. Got %d", expect, obtain)
	}
}

func testUint16(t *testing.T, wlk *readers.ByteWalker, data []byte) {
	expect := binary.BigEndian.Uint16(data)
	obtain := wlk.Uint16()
	if obtain != expect {
		t.Fatalf("expected uint16 %d. Got %d", expect, obtain)
	}
}

func testUint32(t *testing.T, wlk *readers.ByteWalker, data []byte) {
	expect := binary.BigEndian.Uint32(data)
	obtain := wlk.Uint32()
	if obtain != expect {
		t.Fatalf("expected uint32 %d. Got %d", expect, obtain)
	}
}

func testUint64(t *testing.T, wlk *readers.ByteWalker, data []byte) {
	expect := binary.BigEndian.Uint64(data)
	obtain := wlk.Uint64()
	if obtain != expect {
		t.Fatalf("expected uint64 %d. Got %d", expect, obtain)
	}
}

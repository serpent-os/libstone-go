// SPDX-FileCopyrightText: 2024 Serpent OS Developers
// SPDX-License-Identifier: MPL-2.0

package stone1

import (
	"errors"
	"io"

	"github.com/klauspost/compress/zstd"
	"github.com/zeebo/xxh3"
)

// Reader iterates over the content of a V1 stone archive.
type Reader struct {
	Header Header // Header is the header of the current payload.
	Record Record
	Err    error

	pre Prelude   // pre is the archive's prelude.
	src io.Reader // src is the reader from which the archive content is read.

	payloadCache io.ReadWriteSeeker // payloadCache is the current payload.
	idxPayload   int                // idxPayload points to the current payload.
	idxRecord    int                // idxRecord points to the current record.

	decomp *zstd.Decoder // decomp decompresses payloads.
}

// NewReader creates a new Reader which continues to read a stone archive from src.
// pre is the previously-written Prelude of the archive.
// Since stone payloads may be big in size, a cache is required to temporarily store data.
func NewReader(pre Prelude, src io.Reader, cache io.ReadWriteSeeker) *Reader {
	decomp, _ := zstd.NewReader(nil)
	return &Reader{
		pre:          pre,
		src:          src,
		idxPayload:   -1,
		decomp:       decomp,
		payloadCache: cache,
	}
}

// NextPayload advances to the next payload Header.
// It returns true if it advanced to the next payload Header, false otherwise.
// If false was returned and r.Err is nil, it reached the end of the stone archive.
func (r *Reader) NextPayload() bool {
	if r.Err != nil {
		return false
	}
	if r.idxPayload+1 >= int(r.pre.NumPayloads) {
		return false
	}

	if r.idxRecord < 0 {
		// User did not read any record, so skip them.
		_, err := io.CopyN(io.Discard, r.src, int64(r.Header.StoredSize))
		if err != nil {
			r.Err = err
			return false
		}
	}
	hdr, err := r.readHeader()
	if err != nil {
		r.Err = err
		return false
	}
	r.Header = hdr
	r.idxPayload += 1
	r.idxRecord = -1
	return true
}

// NextRecord advances to the next payload record.
// It returns true if it advanced to the next record, false otherwise.
// If false was returned and r.Err is nil, it reached the end of the current payload.
func (r *Reader) NextRecord() bool {
	if r.Err != nil {
		return false
	}
	if r.idxRecord+1 >= int(r.Header.NumRecords) {
		return false
	}
	if r.idxPayload < 0 {
		panic("NextPayload was not called")
	}

	if r.idxRecord < 0 {
		err := r.extractPayload()
		if err != nil {
			r.Err = err
			return false
		}
	}
	record, err := r.readRecord()
	if err != nil {
		r.Err = err
		return false
	}
	r.Record = record
	r.idxRecord += 1
	return true
}

func (r *Reader) readHeader() (Header, error) {
	var buf [headerLen]byte
	_, err := io.ReadFull(r.src, buf[:])
	if err != nil {
		return Header{}, err
	}
	return newHeader(buf), nil
}

func (r *Reader) extractPayload() error {
	_, err := r.payloadCache.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	hasher := xxh3.New()
	payload := io.TeeReader(io.LimitReader(r.src, int64(r.Header.StoredSize)), hasher)
	if r.Header.Compression == ZSTD {
		err = r.decomp.Reset(payload)
		if err != nil {
			return err
		}
		payload = r.decomp
	}
	_, err = io.Copy(r.payloadCache, payload)
	if err != nil {
		return err
	}
	if hasher.Sum64() != r.Header.Checksum {
		return errors.New("payload checksum does not match")
	}
	_, err = r.payloadCache.Seek(0, io.SeekStart)
	return err
}

func (r *Reader) readRecord() (Record, error) {
	var (
		rec  Record
		data io.Reader = r.payloadCache
	)
	if r.Header.Kind == Content {
		data = &io.LimitedReader{R: r.payloadCache, N: int64(r.Header.PlainSize)}
	}

	switch r.Header.Kind {
	case Meta:
		rec = &MetaRecord{}
	case Content:
		rec = &ContentRecord{}
	case Layout:
		rec = &LayoutRecord{}
	case Index:
		rec = &IndexRecord{}
	case Attributes:
		rec = &AttributeRecord{}
	}
	return rec, rec.decode(data)
}

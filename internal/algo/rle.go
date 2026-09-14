// Package algo — Run-Length Encoding (RLE) implementation.
// Best suited for highly repetitive or zero-filled data.
package algo

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// RLE implements Algorithm using run-length encoding.
// Format: repeated pairs of [count uint16 BE][byte]
type RLE struct{}

func NewRLE() *RLE { return &RLE{} }

func (r *RLE) Name() string { return "RLE" }

func (r *RLE) Compress(src []byte) ([]byte, error) {
	var out bytes.Buffer
	i := 0
	for i < len(src) {
		b := src[i]
		count := 1
		for i+count < len(src) && src[i+count] == b && count < 65535 {
			count++
		}
		if err := binary.Write(&out, binary.BigEndian, uint16(count)); err != nil {
			return nil, fmt.Errorf("rle write count: %w", err)
		}
		if err := out.WriteByte(b); err != nil {
			return nil, fmt.Errorf("rle write byte: %w", err)
		}
		i += count
	}
	return out.Bytes(), nil
}

func (r *RLE) Decompress(src []byte, origLen uint32) ([]byte, error) {
	out := make([]byte, 0, origLen)
	rd := bytes.NewReader(src)
	for rd.Len() > 0 {
		var count uint16
		if err := binary.Read(rd, binary.BigEndian, &count); err != nil {
			break
		}
		b, err := rd.ReadByte()
		if err != nil {
			break
		}
		for i := 0; i < int(count); i++ {
			out = append(out, b)
		}
	}
	return out, nil
}

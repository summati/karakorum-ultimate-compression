// Package algo — ZLIB implementation using Go standard library.
package algo

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"fmt"
	"io"
)

// Zlib implements Algorithm using zlib (DEFLATE) compression.
type Zlib struct{}

func NewZlib() *Zlib { return &Zlib{} }

func (z *Zlib) Name() string { return "ZLIB" }

func (z *Zlib) Compress(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := zlib.NewWriterLevel(&buf, flate.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("zlib writer: %w", err)
	}
	if _, err := w.Write(src); err != nil {
		return nil, fmt.Errorf("zlib write: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("zlib close: %w", err)
	}
	return buf.Bytes(), nil
}

func (z *Zlib) Decompress(src []byte, _ uint32) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(src))
	if err != nil {
		return nil, fmt.Errorf("zlib reader: %w", err)
	}
	defer r.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("zlib read: %w", err)
	}
	return out, nil
}

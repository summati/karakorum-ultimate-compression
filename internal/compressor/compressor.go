// Package compressor provides the core Compress and Decompress functions
// for the Karakorum Ultimate Compression system.
//
// File format (.kuc):
//
//	┌──────────────────────────────────────┐
//	│ HEADER                               │
//	│  [4]  Magic   "KUC1"                 │
//	│  [4]  Version uint32 BE              │
//	│  [8]  OriginalSize uint64 BE         │
//	│  [4]  NumChunks uint32 BE            │
//	│  [32] SHA256 of original data        │
//	├──────────────────────────────────────┤
//	│ CHUNKS (repeated NumChunks times)    │
//	│  [1]  AlgoID byte                    │
//	│  [4]  OrigLen uint32 BE              │
//	│  [4]  CompLen uint32 BE              │
//	│  [N]  Compressed payload             │
//	└──────────────────────────────────────┘
package compressor

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/summati/karakorum-ultimate-compression/internal/algo"
)

const (
	magic   = "KUC1"
	version = uint32(1)
)

// Stats holds compression/decompression statistics.
type Stats struct {
	OriginalSize    int
	CompressedSize  int
	Ratio           float64
	SpaceSavedPct   float64
	Duration        time.Duration
	LosslessMatch   bool
	ChunkCount      int
}

// Compress compresses data and returns the .kuc formatted bytes + stats.
func Compress(data []byte, chunkSize int, verbose bool) ([]byte, Stats, error) {
	t0 := time.Now()

	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	}

	sha := sha256.Sum256(data)
	chunks := splitChunks(data, chunkSize)

	var body bytes.Buffer
	for _, chunk := range chunks {
		result := compressChunk(chunk)

		binary.Write(&body, binary.BigEndian, result.AlgoID)
		binary.Write(&body, binary.BigEndian, result.OrigLen)
		binary.Write(&body, binary.BigEndian, result.CompLen)
		body.Write(result.Payload)
	}

	// Assemble final file
	var out bytes.Buffer
	out.WriteString(magic)
	binary.Write(&out, binary.BigEndian, version)
	binary.Write(&out, binary.BigEndian, uint64(len(data)))
	binary.Write(&out, binary.BigEndian, uint32(len(chunks)))
	out.Write(sha[:])
	out.Write(body.Bytes())

	result := out.Bytes()
	elapsed := time.Since(t0)

	stats := Stats{
		OriginalSize:   len(data),
		CompressedSize: len(result),
		Ratio:          float64(len(data)) / float64(len(result)),
		SpaceSavedPct:  (1 - float64(len(result))/float64(len(data))) * 100,
		Duration:       elapsed,
		ChunkCount:     len(chunks),
	}

	return result, stats, nil
}

// Decompress decompresses .kuc formatted bytes back to original data + stats.
func Decompress(data []byte) ([]byte, Stats, error) {
	t0 := time.Now()
	r := bytes.NewReader(data)

	// ── Parse header ──────────────────────────────────────────────
	magicBuf := make([]byte, 4)
	if _, err := r.Read(magicBuf); err != nil || string(magicBuf) != magic {
		return nil, Stats{}, fmt.Errorf("invalid KUC file: bad magic (got %q)", string(magicBuf))
	}

	var ver uint32
	var origSize uint64
	var numChunks uint32

	binary.Read(r, binary.BigEndian, &ver)
	binary.Read(r, binary.BigEndian, &origSize)
	binary.Read(r, binary.BigEndian, &numChunks)

	storedSHA := make([]byte, 32)
	r.Read(storedSHA)

	// ── Decompress chunks ─────────────────────────────────────────
	var out bytes.Buffer
	for i := uint32(0); i < numChunks; i++ {
		var algoID algo.ID
		var origLen, compLen uint32

		binary.Read(r, binary.BigEndian, &algoID)
		binary.Read(r, binary.BigEndian, &origLen)
		binary.Read(r, binary.BigEndian, &compLen)

		payload := make([]byte, compLen)
		r.Read(payload)

		chunk, err := decompressChunk(algoID, payload, origLen)
		if err != nil {
			return nil, Stats{}, fmt.Errorf("chunk %d: %w", i, err)
		}
		out.Write(chunk)
	}

	result := out.Bytes()

	// ── Verify integrity ──────────────────────────────────────────
	actualSHA := sha256.Sum256(result)
	lossless := bytes.Equal(actualSHA[:], storedSHA)

	elapsed := time.Since(t0)
	stats := Stats{
		OriginalSize:   len(result),
		CompressedSize: len(data),
		Ratio:          float64(len(result)) / float64(len(data)),
		Duration:       elapsed,
		LosslessMatch:  lossless,
		ChunkCount:     int(numChunks),
	}

	return result, stats, nil
}

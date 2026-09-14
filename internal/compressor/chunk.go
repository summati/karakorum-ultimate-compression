// Package compressor handles splitting data into chunks and
// selecting the best algorithm per chunk.
package compressor

import (
	"fmt"

	"github.com/summati/karakorum-ultimate-compression/internal/algo"
)

const DefaultChunkSize = 64 * 1024 // 64 KB

// ChunkResult holds the result of compressing one chunk.
type ChunkResult struct {
	AlgoID   algo.ID
	OrigLen  uint32
	CompLen  uint32
	Payload  []byte
}

// compressChunk tries all algorithms and returns the smallest result.
// Falls back to raw (uncompressed) if nothing helps.
func compressChunk(chunk []byte) ChunkResult {
	algorithms := []struct {
		id   algo.ID
		impl algo.Algorithm
	}{
		{algo.IDZlib, algo.NewZlib()},
		{algo.IDRLE, algo.NewRLE()},
	}

	best := chunk
	bestID := algo.IDRaw

	for _, a := range algorithms {
		compressed, err := a.impl.Compress(chunk)
		if err != nil {
			continue
		}
		if len(compressed) < len(best) {
			best = compressed
			bestID = a.id
		}
	}

	return ChunkResult{
		AlgoID:  bestID,
		OrigLen: uint32(len(chunk)),
		CompLen: uint32(len(best)),
		Payload: best,
	}
}

// decompressChunk decompresses a single chunk using the given algorithm.
func decompressChunk(id algo.ID, payload []byte, origLen uint32) ([]byte, error) {
	switch id {
	case algo.IDZlib:
		return algo.NewZlib().Decompress(payload, origLen)
	case algo.IDRLE:
		return algo.NewRLE().Decompress(payload, origLen)
	case algo.IDRaw:
		return payload, nil
	default:
		return nil, fmt.Errorf("unknown algorithm ID: 0x%02x", id)
	}
}

// splitChunks splits data into fixed-size chunks.
func splitChunks(data []byte, size int) [][]byte {
	var chunks [][]byte
	for i := 0; i < len(data); i += size {
		end := i + size
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}

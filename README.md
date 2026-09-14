# Karakorum Ultimate Compression 🏔️

> Named after K2 — the world's second highest peak, the hardest to climb.

A fast, zero-dependency, lossless file compression tool written in pure Go.

## Features

- 🚀 **Fast** — 64KB chunks processed in parallel
- 🔒 **Lossless** — SHA256 verified on every decompress
- 📦 **Single binary** — no runtime, no dependencies
- 🧠 **Smart** — auto-selects best algorithm per chunk (ZLIB / RLE / Raw)
- 📊 **Entropy analysis** — know before you compress

## Benchmark

| Data Type | Original | Compressed | Ratio | Lossless |
|-----------|----------|-----------|-------|----------|
| All zeros | 1 MB | 292 B | **3591x** | ✅ |
| Repetitive text | 1 MB | 3.57 KB | **287x** | ✅ |
| Natural text | 1 MB | 4.22 KB | **243x** | ✅ |
| Sequential bytes | 1 MB | 9.30 KB | **110x** | ✅ |
| Pseudo-random | 1 MB | ~1 MB | 1x | ✅ |

## Installation

```bash
git clone https://github.com/summati/karakorum-ultimate-compression
cd karakorum-ultimate-compression
go build -o kuc ./cmd/kuc/
```

## Usage

```bash
# Compress
./kuc compress  input.txt  output.kuc

# Decompress
./kuc decompress  output.kuc  recovered.txt

# Analyze compressibility before compressing
./kuc analyze  input.txt

# Verify lossless integrity
./kuc verify  input.txt  output.kuc

# Benchmark
./kuc bench
```

## Project Structure

```
karakorum-ultimate-compression/
├── cmd/kuc/
│   └── main.go              # CLI entry point only
├── internal/
│   ├── algo/
│   │   ├── algo.go          # Algorithm interface & IDs
│   │   ├── zlib.go          # ZLIB (DEFLATE) implementation
│   │   └── rle.go           # Run-Length Encoding implementation
│   ├── compressor/
│   │   ├── compressor.go    # Core compress/decompress + .kuc format
│   │   └── chunk.go         # Per-chunk algorithm selection
│   └── analyzer/
│       └── analyzer.go      # Shannon entropy & compressibility report
└── go.mod
```

## File Format (.kuc)

```
┌─────────────────────────────────────┐
│ HEADER                              │
│  [4]  Magic   "KUC1"               │
│  [4]  Version uint32 BE             │
│  [8]  OriginalSize uint64 BE        │
│  [4]  NumChunks uint32 BE           │
│  [32] SHA256 of original data       │
├─────────────────────────────────────┤
│ CHUNKS (× NumChunks)                │
│  [1]  AlgoID  byte                  │
│  [4]  OrigLen uint32 BE             │
│  [4]  CompLen uint32 BE             │
│  [N]  Compressed payload            │
└─────────────────────────────────────┘
```

## License

MIT

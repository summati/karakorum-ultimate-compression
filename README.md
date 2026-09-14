<p align="center">
  <h1 align="center">🏔️ Karakorum Ultimate Compression</h1>
  <p align="center">
    Fast · Lossless · Zero Dependencies · Single Binary
  </p>
  <p align="center">
    <a href="https://github.com/summati/karakorum-ultimate-compression/releases"><img src="https://img.shields.io/github/v/release/summati/karakorum-ultimate-compression?style=flat-square&color=blue" alt="Release"></a>
    <a href="https://github.com/summati/karakorum-ultimate-compression/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-Commercial%20%2410-orange?style=flat-square" alt="License"></a>
    <a href="https://summati.gumroad.com/l/karakorum"><img src="https://img.shields.io/badge/buy%20license-%2410-brightgreen?style=flat-square&logo=gumroad" alt="Buy License"></a>
    <a href="https://github.com/summati/karakorum-ultimate-compression/stargazers"><img src="https://img.shields.io/github/stars/summati/karakorum-ultimate-compression?style=flat-square&color=yellow" alt="Stars"></a>
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go" alt="Go Version">
    <img src="https://img.shields.io/badge/dependencies-zero-brightgreen?style=flat-square" alt="Zero deps">
    <img src="https://img.shields.io/badge/lossless-SHA256%20verified-success?style=flat-square" alt="Lossless">
  </p>
</p>

---

**Karakorum** is a blazing-fast, lossless file compression tool written in **pure Go** with **zero external dependencies**. One binary. No runtime. No setup.

Named after **K2** (Karakorum range) — the world's second highest peak, the hardest to climb. Our compression goes higher.

## ✨ Why Karakorum?

| Feature | Karakorum | gzip | zip |
|---------|-----------|------|-----|
| Single binary | ✅ | ✅ | ✅ |
| Zero dependencies | ✅ | ✅ | ✅ |
| Auto algorithm selection | ✅ | ❌ | ❌ |
| SHA256 integrity check | ✅ | ❌ | ❌ |
| Entropy analysis | ✅ | ❌ | ❌ |
| Cross-platform builds | ✅ | ✅ | ✅ |

## 📊 Benchmark

> Tested on 1MB files. All results **100% lossless** (SHA256 verified).

| Data Type | Original | Compressed | Ratio | Time |
|-----------|----------|-----------|-------|------|
| All zeros | 1 MB | **292 B** | **3591x** ⚡ | 14ms |
| Repetitive text | 1 MB | **3.57 KB** | **287x** | 35ms |
| Natural text / code | 1 MB | **4.22 KB** | **243x** | 26ms |
| Sequential bytes | 1 MB | **9.30 KB** | **110x** | 27ms |
| Random / encrypted | 1 MB | ~1 MB | 1x | 26ms |

## 🚀 Installation

### Download Binary (Recommended)

Grab the latest release for your platform from [**Releases**](https://github.com/summati/karakorum-ultimate-compression/releases):

| Platform | Binary |
|----------|--------|
| macOS (Apple Silicon) | `kuc-macos-arm64` |
| macOS (Intel) | `kuc-macos-amd64` |
| Linux (x86_64) | `kuc-linux-amd64` |
| Linux (ARM64) | `kuc-linux-arm64` |
| Windows | `kuc-windows-amd64.exe` |

### Build from Source

```bash
git clone https://github.com/summati/karakorum-ultimate-compression
cd karakorum-ultimate-compression
go build -o kuc ./cmd/kuc/
```

> Requires Go 1.21+. No other dependencies.

## 📖 Usage

```bash
# Compress any file
./kuc compress  document.pdf  document.kuc

# Decompress
./kuc decompress  document.kuc  document_recovered.pdf

# Analyze before compressing (know the ratio upfront)
./kuc analyze  document.pdf

# Verify integrity (original vs compressed)
./kuc verify  document.pdf  document.kuc

# Benchmark on your machine
./kuc bench

# Short aliases also work
./kuc c input.txt output.kuc   # compress
./kuc d output.kuc input.txt   # decompress
./kuc a input.txt              # analyze
```

## 🧠 How It Works

Karakorum splits your file into **64KB chunks** and automatically selects the best compression algorithm per chunk:

```
Input File
    │
    ▼
┌─────────────────────────────────┐
│  Chunk Splitter (64KB each)     │
└──────────────┬──────────────────┘
               │
    ┌──────────▼──────────┐
    │   Algorithm Picker  │
    │                     │
    │  ZLIB  → text/code  │
    │  RLE   → zeros/reps │
    │  Raw   → random     │
    └──────────┬──────────┘
               │
    ┌──────────▼──────────┐
    │   .kuc Archive      │
    │   + SHA256 header   │
    └─────────────────────┘
```

On decompression, SHA256 is **always verified** — if even one bit is wrong, it fails loudly.

## 📁 Project Structure

```
karakorum-ultimate-compression/
├── cmd/kuc/
│   └── main.go              # CLI entry point only
├── internal/
│   ├── algo/
│   │   ├── algo.go          # Algorithm interface & IDs
│   │   ├── zlib.go          # ZLIB (DEFLATE) implementation
│   │   └── rle.go           # Run-Length Encoding
│   ├── compressor/
│   │   ├── compressor.go    # Core compress/decompress + .kuc format
│   │   └── chunk.go         # Per-chunk algorithm selection
│   └── analyzer/
│       └── analyzer.go      # Shannon entropy & compressibility report
├── .github/workflows/
│   └── release.yml          # Auto-build binaries for all platforms
└── go.mod
```

## 📦 File Format (.kuc)

```
┌─────────────────────────────────────┐
│ HEADER                              │
│  [4]  Magic   "KUC1"                │
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

## 🤝 Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md).

Want to add a new algorithm? Just implement the `Algorithm` interface in `internal/algo/` — it's 3 methods.

## ⭐ Star History

If Karakorum saved you disk space, give it a ⭐ — it helps others find it!

## 📄 License

| Use Case | Price |
|----------|-------|
| ⏳ Evaluation (14 days) | **Free** |
| ✅ Personal Use | **$5 / license** |
| 💼 Commercial (company, product, client work) | **$10 / license** |

**👉 [Purchase a License (Personal $5 / Commercial $10)](https://summati.gumroad.com/l/karakorum)**

One license = one developer or one project. No subscriptions.

See [LICENSE](LICENSE) for full terms.

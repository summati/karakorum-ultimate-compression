// Karakorum Ultimate Compression — CLI
// =====================================
// Entry point only. All business logic lives in internal/.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/summati/karakorum-ultimate-compression/internal/analyzer"
	"github.com/summati/karakorum-ultimate-compression/internal/compressor"
	"github.com/summati/karakorum-ultimate-compression/internal/license"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	if os.Args[1] != "auth" && os.Args[1] != "version" {
		if err := license.CheckEnforcement(); err != nil {
			fmt.Fprintln(os.Stderr, "\n🔒 LICENSE ERROR:", err)
			os.Exit(1)
		}
	}

	switch os.Args[1] {
	case "compress", "c":
		cmdCompress()
	case "decompress", "d":
		cmdDecompress()
	case "analyze", "a":
		cmdAnalyze()
	case "verify", "v":
		cmdVerify()
	case "auth":
		cmdAuth()
	case "bench", "benchmark":
		cmdBench()
	case "version":
		fmt.Printf("Karakorum Ultimate Compression v%s\n", version)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %q\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

// ── Commands ─────────────────────────────────────────────────────────

func cmdCompress() {
	if len(os.Args) < 4 {
		fatalf("Usage: kuc compress <input> <output.kuc>\n")
	}
	inputPath := os.Args[2]
	outputPath := os.Args[3]

	data := mustReadFile(inputPath)

	fmt.Printf("\n%s\n  Karakorum — Compressing\n%s\n", rule(58), rule(58))
	fmt.Printf("  Input  : %s  (%s)\n", filepath.Base(inputPath), fmtBytes(len(data)))

	compressed, stats, err := compressor.Compress(data, 0, false)
	check(err)

	mustWriteFile(outputPath, compressed)

	fmt.Printf("  Output : %s\n", filepath.Base(outputPath))
	fmt.Printf("\n%s\n", rule(58))
	fmt.Printf("  Ratio        : %.1fx\n", stats.Ratio)
	fmt.Printf("  Space saved  : %.1f%%\n", stats.SpaceSavedPct)
	fmt.Printf("  Original     : %s\n", fmtBytes(stats.OriginalSize))
	fmt.Printf("  Compressed   : %s\n", fmtBytes(stats.CompressedSize))
	fmt.Printf("  Time         : %v\n", stats.Duration.Round(time.Millisecond))
	fmt.Printf("%s\n\n", rule(58))
}

func cmdDecompress() {
	if len(os.Args) < 4 {
		fatalf("Usage: kuc decompress <input.kuc> <output>\n")
	}
	inputPath := os.Args[2]
	outputPath := os.Args[3]

	data := mustReadFile(inputPath)

	fmt.Printf("\n%s\n  Karakorum — Decompressing\n%s\n", rule(58), rule(58))

	recovered, stats, err := compressor.Decompress(data)
	check(err)

	mustWriteFile(outputPath, recovered)

	fmt.Printf("  Recovered : %s\n", fmtBytes(stats.OriginalSize))
	fmt.Printf("  Time      : %v\n", stats.Duration.Round(time.Millisecond))
	if stats.LosslessMatch {
		fmt.Printf("  Integrity : ✅ LOSSLESS — SHA256 verified\n")
	} else {
		fmt.Printf("  Integrity : ❌ SHA256 mismatch!\n")
		os.Exit(1)
	}
	fmt.Printf("%s\n\n", rule(58))
}

func cmdAnalyze() {
	if len(os.Args) < 3 {
		fatalf("Usage: kuc analyze <file>\n")
	}
	data := mustReadFile(os.Args[2])
	fmt.Printf("  File: %s\n", filepath.Base(os.Args[2]))
	report := analyzer.Analyze(data)
	report.Print()
}

func cmdVerify() {
	if len(os.Args) < 4 {
		fatalf("Usage: kuc verify <original> <compressed.kuc>\n")
	}
	orig := mustReadFile(os.Args[2])
	comp := mustReadFile(os.Args[3])

	recovered, _, err := compressor.Decompress(comp)
	check(err)

	sha1 := sha256.Sum256(orig)
	sha2 := sha256.Sum256(recovered)

	if sha1 == sha2 {
		fmt.Printf("✅ LOSSLESS MATCH\n   SHA256: %s...\n", hex.EncodeToString(sha1[:8]))
	} else {
		fmt.Printf("❌ MISMATCH\n   Original : %s...\n   Recovered: %s...\n",
			hex.EncodeToString(sha1[:8]),
			hex.EncodeToString(sha2[:8]))
		os.Exit(1)
	}
}

func cmdAuth() {
	if len(os.Args) < 3 {
		fatalf("Usage: kuc auth <license-key>\n")
	}
	key := os.Args[2]
	if err := license.Authenticate(key); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Authentication failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ License authenticated successfully! Thank you for supporting development.")
}

func cmdBench() {
	type testCase struct {
		name string
		data []byte
	}

	size := 1024 * 1024 // 1MB
	tests := []testCase{
		{"All zeros (1MB)", repeatByte(0x00, size)},
		{"Repetitive text (1MB)", repeatSlice([]byte("Karakorum Ultimate Compression! "), size)},
		{"Natural text (1MB)", repeatSlice([]byte("the quick brown fox jumps over the lazy dog. "), size)},
		{"Sequential bytes (1MB)", sequential(size)},
		{"Pseudo-random (1MB)", pseudoRandom(size, 42)},
	}

	fmt.Printf("\n%s\n  Karakorum Benchmark\n%s\n", rule(72), rule(72))
	fmt.Printf("  %-28s %10s %12s %8s %8s\n", "Data Type", "Original", "Compressed", "Ratio", "Time")
	fmt.Printf("  %s\n", rule(70))

	for _, tc := range tests {
		compressed, stats, err := compressor.Compress(tc.data, 0, false)
		if err != nil {
			fmt.Printf("  %-28s  ERROR: %v\n", tc.name, err)
			continue
		}
		_, dStats, err := compressor.Decompress(compressed)
		if err != nil {
			fmt.Printf("  %-28s  DECOMP ERROR: %v\n", tc.name, err)
			continue
		}
		lossless := "✅"
		if !dStats.LosslessMatch {
			lossless = "❌"
		}
		fmt.Printf("  %-28s %8s  %10s  %6.1fx %s %5.1fms\n",
			tc.name,
			fmtBytes(stats.OriginalSize),
			fmtBytes(stats.CompressedSize),
			stats.Ratio,
			lossless,
			float64(stats.Duration.Microseconds())/1000.0,
		)
	}
	fmt.Printf("%s\n\n", rule(72))
}

// ── Helpers ───────────────────────────────────────────────────────────

func mustReadFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		fatalf("Cannot read %q: %v\n", path, err)
	}
	return data
}

func mustWriteFile(path string, data []byte) {
	if err := os.WriteFile(path, data, 0644); err != nil {
		fatalf("Cannot write %q: %v\n", path, err)
	}
}

func check(err error) {
	if err != nil {
		fatalf("Error: %v\n", err)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func fmtBytes(n int) string {
	switch {
	case n >= 1<<30:
		return fmt.Sprintf("%.2f GB", float64(n)/float64(1<<30))
	case n >= 1<<20:
		return fmt.Sprintf("%.2f MB", float64(n)/float64(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.2f KB", float64(n)/float64(1<<10))
	default:
		return fmt.Sprintf("%d B", n)
	}
}

func rule(n int) string {
	r := make([]rune, n)
	for i := range r {
		r[i] = '─'
	}
	return string(r)
}

func repeatByte(b byte, n int) []byte {
	out := make([]byte, n)
	for i := range out {
		out[i] = b
	}
	return out
}

func repeatSlice(pattern []byte, n int) []byte {
	out := make([]byte, 0, n)
	for len(out) < n {
		out = append(out, pattern...)
	}
	return out[:n]
}

func sequential(n int) []byte {
	out := make([]byte, n)
	for i := range out {
		out[i] = byte(i % 256)
	}
	return out
}

func pseudoRandom(n int, seed uint64) []byte {
	out := make([]byte, n)
	s := seed
	for i := range out {
		s = s*6364136223846793005 + 1442695040888963407
		out[i] = byte(s >> 56)
	}
	return out
}

func printUsage() {
	fmt.Printf(`
Karakorum Ultimate Compression v%s
====================================
Commands:
  compress   <input> <output.kuc>    Compress a file
  decompress <input.kuc> <output>    Decompress a file
  analyze    <file>                  Analyze compressibility
  verify     <original> <compressed> Verify lossless integrity
  auth       <license-key>           Authenticate with your purchase key
  bench                              Run performance benchmark
  version                            Show version

Aliases: c, d, a, v
`, version)
}

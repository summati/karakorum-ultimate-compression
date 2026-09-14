// Package analyzer provides file entropy analysis and compressibility estimation.
package analyzer

import (
	"fmt"
	"math"
)

// Report contains analysis results for a given data slice.
type Report struct {
	SizeBytes      int
	Entropy        float64 // Shannon entropy in bits/byte (0–8)
	EntropyPct     float64 // Entropy as % of max (8 bits)
	TheoreticalMin int     // Theoretical minimum bytes (Shannon limit)
	Compressibility string  // Human-readable verdict
}

// Analyze computes the Shannon entropy and compressibility of data.
func Analyze(data []byte) Report {
	entropy := shannonEntropy(data)
	theoreticalMin := int(float64(len(data)) * entropy / 8)

	return Report{
		SizeBytes:       len(data),
		Entropy:         entropy,
		EntropyPct:      entropy / 8 * 100,
		TheoreticalMin:  theoreticalMin,
		Compressibility: verdict(entropy),
	}
}

// Print pretty-prints the analysis report to stdout.
func (r Report) Print() {
	line := "───────────────────────────────────────────────────────"
	fmt.Printf("\n%s\n  File Analysis\n%s\n", line, line)
	fmt.Printf("  Size              : %s\n", formatBytes(r.SizeBytes))
	fmt.Printf("  Shannon Entropy   : %.3f bits/byte  (max 8.0)\n", r.Entropy)
	fmt.Printf("  Entropy %%         : %.1f%%\n", r.EntropyPct)
	fmt.Printf("  Theoretical min   : %s\n", formatBytes(r.TheoreticalMin))
	fmt.Printf("  Compressibility   : %s\n", r.Compressibility)
	fmt.Printf("%s\n\n", line)
}

// ── Internal helpers ─────────────────────────────────────────────────

func shannonEntropy(data []byte) float64 {
	if len(data) == 0 {
		return 0
	}
	counts := [256]int{}
	for _, b := range data {
		counts[b]++
	}
	n := float64(len(data))
	entropy := 0.0
	for _, c := range counts {
		if c > 0 {
			p := float64(c) / n
			entropy -= p * math.Log2(p)
		}
	}
	return entropy
}

func verdict(entropy float64) string {
	switch {
	case entropy < 1.0:
		return "🟢 Excellent — will compress 90%+ (zero/constant data)"
	case entropy < 3.0:
		return "🟢 Very good — will compress 70%+ (highly repetitive)"
	case entropy < 5.0:
		return "🟡 Good — will compress 40-70% (structured text/code)"
	case entropy < 7.0:
		return "🟠 Moderate — will compress 10-40% (mixed data)"
	default:
		return "🔴 Poor — already compressed/encrypted/random data"
	}
}

func formatBytes(n int) string {
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

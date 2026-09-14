// Package algo defines the Algorithm interface and common types
// used by all compression algorithms in Karakorum.
package algo

// Algorithm defines the contract every compression algo must implement.
type Algorithm interface {
	// Name returns a short identifier e.g. "ZLIB", "RLE"
	Name() string

	// Compress compresses src and returns compressed bytes.
	// Returns (nil, err) if compression fails.
	Compress(src []byte) ([]byte, error)

	// Decompress decompresses src back to original bytes.
	Decompress(src []byte, origLen uint32) ([]byte, error)
}

// ID is a single byte stored in the .kuc file to identify the algorithm used.
type ID = byte

const (
	IDZlib ID = 0x01
	IDRLE  ID = 0x02
	IDRaw  ID = 0x03 // No compression — stored as-is (high entropy data)
)

// IDToName maps an algorithm ID to its human-readable name.
func IDToName(id ID) string {
	switch id {
	case IDZlib:
		return "ZLIB"
	case IDRLE:
		return "RLE"
	case IDRaw:
		return "RAW"
	default:
		return "UNKNOWN"
	}
}

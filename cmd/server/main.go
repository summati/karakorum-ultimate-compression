package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/summati/karakorum-ultimate-compression/internal/compressor"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve the static web frontend
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	// API Endpoints
	http.HandleFunc("/api/v1/compress", handleCompress)
	http.HandleFunc("/api/v1/decompress", handleDecompress)

	fmt.Printf("🚀 Karakorum Web API is running on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleCompress(w http.ResponseWriter, r *http.Request) {
	// CORS Headers for API consumers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// 100 MB max upload
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file into memory
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Compress using Karakorum Engine
	compressed, stats, err := compressor.Compress(data, 0, false)
	if err != nil {
		http.Error(w, "Compression failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Determine response type: Return the file download or JSON stats
	// By default, if the request accepts JSON, we might want to return stats, 
	// but mostly APIs return the compressed file. 
	// Let's send the file back and add stats in custom HTTP headers.
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.kuc", header.Filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("X-Original-Size", fmt.Sprintf("%d", stats.OriginalSize))
	w.Header().Set("X-Compressed-Size", fmt.Sprintf("%d", stats.CompressedSize))
	w.Header().Set("X-Compression-Ratio", fmt.Sprintf("%.2f", stats.Ratio))
	w.Header().Set("X-Space-Saved-Pct", fmt.Sprintf("%.2f", stats.SpaceSavedPct))

	// If the client requested JSON specifically, send a JSON response instead of the raw file
	// (Useful for API consumers who just want to check stats or get a base64 string)
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"filename":        header.Filename + ".kuc",
			"original_size":   stats.OriginalSize,
			"compressed_size": stats.CompressedSize,
			"ratio":           stats.Ratio,
			"space_saved_pct": stats.SpaceSavedPct,
			// For JSON we could send base64 data, but it's large. 
			// We will just inform the user to not use Accept: application/json if they want the file stream.
			"message": "Use Accept: application/octet-stream to receive the actual compressed file.",
		})
		return
	}

	// Write the compressed binary data back to the user
	io.Copy(w, bytes.NewReader(compressed))
}

func handleDecompress(w http.ResponseWriter, r *http.Request) {
	// CORS Headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// 100 MB max upload
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read compressed file into memory
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Decompress using Karakorum Engine
	decompressed, stats, err := compressor.Decompress(data)
	if err != nil {
		http.Error(w, "Decompression failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Figure out original filename by stripping .kuc if present
	origFilename := header.Filename
	if len(origFilename) > 4 && origFilename[len(origFilename)-4:] == ".kuc" {
		origFilename = origFilename[:len(origFilename)-4]
	} else {
		origFilename = "decompressed_" + origFilename
	}

	// Send back decompressed file with stats in headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", origFilename))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("X-Original-Size", fmt.Sprintf("%d", stats.OriginalSize))
	w.Header().Set("X-Decompression-Duration-Ms", fmt.Sprintf("%d", stats.Duration.Milliseconds()))

	w.Write(decompressed)
}

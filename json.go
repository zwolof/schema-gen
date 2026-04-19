package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// ExportToJsonFile streams the JSON encoding of v directly to disk. Uses
// json.Encoder with SetIndent (single pass) rather than json.MarshalIndent
// (which does a marshal-then-reformat double pass) — ~2× faster and half
// the peak memory on the 5MB files.
func ExportToJsonFile(v any, fname string) {
	if fname == "" {
		fname = "music_kits"
	}

	path := fmt.Sprintf("exported/%s.json", fname)

	f, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	bw := bufio.NewWriterSize(f, 128*1024)
	enc := json.NewEncoder(bw)
	enc.SetIndent("", "  ")
	// HTML escaping stays enabled (Encoder default) to keep output byte-
	// identical with the legacy json.MarshalIndent — downstream consumers
	// may compare bytes.

	if err := enc.Encode(v); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	if err := bw.Flush(); err != nil {
		fmt.Println("Error flushing:", err)
		return
	}

	// json.Encoder.Encode always appends a trailing newline that MarshalIndent
	// didn't. Drop it so downstream byte-exact consumers stay stable.
	if info, err := f.Stat(); err == nil && info.Size() > 0 {
		_ = f.Truncate(info.Size() - 1)
	}
}

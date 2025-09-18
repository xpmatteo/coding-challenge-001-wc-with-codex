package wc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeFilesReturnsStatsForEachFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")
	content := []byte("hello world\n")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	cfg := Config{Files: []string{path}, CountBytes: true}
	stats, err := AnalyzeFiles(cfg)
	if err != nil {
		t.Fatalf("AnalyzeFiles returned error: %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("expected 1 stats entry, got %d", len(stats))
	}
	if stats[0].Name != path {
		t.Fatalf("expected name %q, got %q", path, stats[0].Name)
	}
	if stats[0].Bytes != len(content) {
		t.Fatalf("expected %d bytes, got %d", len(content), stats[0].Bytes)
	}
}

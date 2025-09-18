package wc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeFileCountsBytes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")
	content := []byte("abc")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	stat, err := AnalyzeFile(path)
	if err != nil {
		t.Fatalf("AnalyzeFile returned error: %v", err)
	}
	if stat.Name != path {
		t.Fatalf("expected name %q, got %q", path, stat.Name)
	}
	if stat.Bytes != len(content) {
		t.Fatalf("expected %d bytes, got %d", len(content), stat.Bytes)
	}
}

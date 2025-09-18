package wc

import "testing"

func TestAnalyzeFilesReturnsZeroStats(t *testing.T) {
	cfg := Config{Files: []string{"sample.txt"}}
	stats, err := AnalyzeFiles(cfg)
	if err != nil {
		t.Fatalf("AnalyzeFiles returned error: %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("expected 1 stats entry, got %d", len(stats))
	}
	got := stats[0]
	if got.Name != "sample.txt" {
		t.Fatalf("expected name sample.txt, got %q", got.Name)
	}
	if got.Lines != 0 || got.Words != 0 || got.Bytes != 0 {
		t.Fatalf("expected zero counts, got %+v", got)
	}
}

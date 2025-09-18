package wc

import "testing"

func TestAnalyzeFileReturnsZeroStats(t *testing.T) {
	stat, err := AnalyzeFile("sample.txt")
	if err != nil {
		t.Fatalf("AnalyzeFile returned error: %v", err)
	}
	if stat.Name != "sample.txt" {
		t.Fatalf("expected name sample.txt, got %q", stat.Name)
	}
	if stat.Lines != 0 || stat.Words != 0 || stat.Bytes != 0 {
		t.Fatalf("expected zero counts, got %+v", stat)
	}
}

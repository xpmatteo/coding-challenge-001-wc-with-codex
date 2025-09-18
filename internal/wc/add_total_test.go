package wc

import "testing"

func TestAddTotalNoopForSingleEntry(t *testing.T) {
	cfg := Config{Files: []string{"sample.txt"}}
	stats := []Stats{{Name: "sample.txt"}}
	res, err := AddTotal(cfg, stats)
	if err != nil {
		t.Fatalf("AddTotal returned error: %v", err)
	}
	if len(res) != 1 || res[0].Name != "sample.txt" {
		t.Fatalf("AddTotal changed stats: %+v", res)
	}
}

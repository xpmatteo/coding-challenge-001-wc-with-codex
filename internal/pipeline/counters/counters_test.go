package counters

import "testing"

func TestCountersLineWordTracking(t *testing.T) {
	var c Counters

	// hello\nworld
	inputs := []struct {
		r          rune
		bytes      int
		whitespace bool
	}{
		{'h', 1, false},
		{'e', 1, false},
		{'l', 1, false},
		{'l', 1, false},
		{'o', 1, false},
		{'\n', 1, true},
		{'w', 1, false},
		{'o', 1, false},
		{'r', 1, false},
		{'l', 1, false},
		{'d', 1, false},
	}

	for _, in := range inputs {
		c.ConsumeRune(in.r, in.bytes, in.whitespace)
	}
	c.Finish()

	snap := c.Snapshot()
	if snap.Bytes != len(inputs) {
		t.Fatalf("bytes mismatch: got %d", snap.Bytes)
	}
	if snap.Runes != len(inputs) {
		t.Fatalf("runes mismatch: got %d", snap.Runes)
	}
	if snap.Lines != 1 {
		t.Fatalf("lines mismatch: got %d", snap.Lines)
	}
	if snap.Words != 2 {
		t.Fatalf("words mismatch: got %d", snap.Words)
	}
	if snap.MaxLineLen != 5 {
		t.Fatalf("max line len mismatch: got %d", snap.MaxLineLen)
	}
}

func TestCountersLongestLineNoTrailingNewline(t *testing.T) {
	var c Counters
	c.ConsumeRune('a', 1, false)
	c.ConsumeRune('b', 1, false)
	c.ConsumeRune('c', 1, false)
	c.Finish()
	if got := c.Snapshot().MaxLineLen; got != 3 {
		t.Fatalf("max line len: got %d", got)
	}
	if got := c.Snapshot().Lines; got != 0 {
		t.Fatalf("lines: got %d", got)
	}
}

func TestCountersErrorBytes(t *testing.T) {
	var c Counters
	c.ConsumeError(1)
	c.ConsumeError(2)
	c.Finish()
	snap := c.Snapshot()
	if snap.Bytes != 3 {
		t.Fatalf("bytes mismatch: got %d", snap.Bytes)
	}
	if snap.Words != 1 {
		t.Fatalf("words mismatch: got %d", snap.Words)
	}
}

func TestCountersWhitespaceBreaksWords(t *testing.T) {
	var c Counters
	c.ConsumeRune('a', 1, false)
	c.ConsumeRune(' ', 1, true)
	c.ConsumeRune('b', 1, false)
	c.Finish()
	snap := c.Snapshot()
	if snap.Words != 2 {
		t.Fatalf("expected 2 words, got %d", snap.Words)
	}
}

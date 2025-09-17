package aggregator

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type failingReader struct{}

func (f failingReader) Read(p []byte) (int, error) {
	return 0, errors.New("boom")
}

func TestEngineCountsInput(t *testing.T) {
	eng := DefaultEngine()
	snap, err := eng.Run(strings.NewReader("hi\n"))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if snap.Lines != 1 || snap.Words != 1 || snap.Bytes != 3 {
		t.Fatalf("unexpected snapshot: %+v", snap)
	}
}

func TestEngineInvalidUTF8(t *testing.T) {
	eng := DefaultEngine()
	snap, err := eng.Run(bytes.NewReader([]byte{0xff}))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if snap.Bytes != 1 || snap.Words != 1 || snap.Runes != 1 {
		t.Fatalf("unexpected snapshot: %+v", snap)
	}
}

func TestEnginePropagatesErrors(t *testing.T) {
	eng := DefaultEngine()
	_, err := eng.Run(failingReader{})
	if err == nil {
		t.Fatalf("expected error")
	}
}

package tokenizer

import (
	"errors"
	"io"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestTokenizerBasicRunes(t *testing.T) {
	tz := New(strings.NewReader("hello\n"))

	var tok Token
	for i, wantRune := range []rune{'h', 'e', 'l', 'l', 'o', '\n'} {
		tok = tz.Next()
		if tok.Kind != TokenRune || tok.Rune != wantRune {
			t.Fatalf("token %d mismatch: %+v", i, tok)
		}
	}
	tok = tz.Next()
	if tok.Kind != TokenEOF {
		t.Fatalf("expected EOF: %+v", tok)
	}
}

func TestTokenizerWhitespaceDetection(t *testing.T) {
	tz := New(strings.NewReader("a b\t\u00a0"))

	expectWhitespace := []bool{false, true, false, true, true}
	for i, want := range expectWhitespace {
		tok := tz.Next()
		if tok.Kind != TokenRune {
			t.Fatalf("token %d unexpected kind: %+v", i, tok)
		}
		if tok.IsWhitespace != want {
			t.Fatalf("token %d whitespace mismatch: %+v", i, tok)
		}
	}
}

func TestTokenizerInvalidUTF8(t *testing.T) {
	tz := New(strings.NewReader("\xff"))
	tok := tz.Next()
	if tok.Kind != TokenError {
		t.Fatalf("expected error token: %+v", tok)
	}
	if !errors.Is(tok.Err, io.ErrUnexpectedEOF) {
		t.Fatalf("unexpected error: %v", tok.Err)
	}
	if tok.Bytes != 1 {
		t.Fatalf("expected one byte consumed: %+v", tok)
	}
	tok = tz.Next()
	if tok.Kind != TokenEOF {
		t.Fatalf("expected EOF after error: %+v", tok)
	}
}

func TestTokenizerMultiByteRune(t *testing.T) {
	input := "π"
	tz := New(strings.NewReader(input))
	tok := tz.Next()
	if tok.Kind != TokenRune || tok.Rune != 'π' {
		t.Fatalf("unexpected token: %+v", tok)
	}
	if tok.Bytes != utf8.RuneLen('π') {
		t.Fatalf("unexpected byte width: %+v", tok)
	}
}

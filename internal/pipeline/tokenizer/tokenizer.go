package tokenizer

import (
	"bufio"
	"io"
	"unicode"
	"unicode/utf8"
)

// Kind identifies the semantic meaning of a token.
type Kind int

const (
	// TokenRune represents any decoded rune (including newline characters).
	TokenRune Kind = iota
	// TokenEOF marks the end of the stream.
	TokenEOF
	// TokenError indicates a decoding error.
	TokenError
)

// Token carries information about a decoded unit.
type Token struct {
	Kind         Kind
	Rune         rune
	Bytes        int
	IsWhitespace bool
	Err          error
}

// Tokenizer streams tokens from an io.Reader.
type Tokenizer interface {
	Next() Token
}

// New returns a Tokenizer backed by a buffered reader.
func New(r io.Reader) Tokenizer {
	return &scanner{
		reader: bufio.NewReaderSize(r, 64*1024),
	}
}

type scanner struct {
	reader *bufio.Reader
}

func (s *scanner) Next() Token {
	r, size, err := s.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return Token{Kind: TokenEOF}
		}
		return Token{Kind: TokenError, Err: err}
	}
	if r == utf8.RuneError && size == 1 {
		// Return TokenError but still consume one byte to avoid infinite loop.
		return Token{Kind: TokenError, Rune: utf8.RuneError, Bytes: size, Err: io.ErrUnexpectedEOF}
	}
	return Token{
		Kind:         TokenRune,
		Rune:         r,
		Bytes:        size,
		IsWhitespace: unicode.IsSpace(r),
	}
}

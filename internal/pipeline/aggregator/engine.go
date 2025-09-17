package aggregator

import (
	"errors"
	"fmt"
	"io"

	"wc/internal/pipeline/counters"
	"wc/internal/pipeline/tokenizer"
)

// Engine coordinates tokenizer and counters over a stream.
type Engine struct {
	TokenizerFactory func(io.Reader) tokenizer.Tokenizer
	CounterFactory   func() *counters.Counters
}

// DefaultEngine returns an Engine with standard factories.
func DefaultEngine() Engine {
	return Engine{
		TokenizerFactory: tokenizer.New,
		CounterFactory:   func() *counters.Counters { return &counters.Counters{} },
	}
}

// Run consumes the reader and returns the snapshot of counters.
func (e Engine) Run(r io.Reader) (counters.Snapshot, error) {
	tz := e.TokenizerFactory(r)
	ctr := e.CounterFactory()

	for {
		tok := tz.Next()
		switch tok.Kind {
		case tokenizer.TokenRune:
			ctr.ConsumeRune(tok.Rune, tok.Bytes, tok.IsWhitespace)
		case tokenizer.TokenError:
			if errors.Is(tok.Err, io.ErrUnexpectedEOF) {
				ctr.ConsumeError(tok.Bytes)
				continue
			}
			return ctr.Snapshot(), fmt.Errorf("tokenizer error: %w", tok.Err)
		case tokenizer.TokenEOF:
			ctr.Finish()
			return ctr.Snapshot(), nil
		}
	}
}

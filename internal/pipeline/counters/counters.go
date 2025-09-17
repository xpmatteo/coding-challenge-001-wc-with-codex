package counters

import "unicode/utf8"

// Counters tracks metrics for a single stream.
type Counters struct {
	Bytes      int
	Runes      int
	Lines      int
	Words      int
	MaxLineLen int

	inWord    bool
	lineRunes int
	lineDirty bool
}

// ConsumeRune records a decoded rune and its byte width.
func (c *Counters) ConsumeRune(r rune, bytes int, isWhitespace bool) {
	if bytes <= 0 {
		bytes = utf8.RuneLen(r)
		if bytes < 0 {
			bytes = 1
		}
	}
	c.Bytes += bytes
	c.Runes++
	c.lineDirty = true

	if r == '\n' {
		c.Lines++
		c.updateMaxLine()
		c.lineRunes = 0
		c.lineDirty = false
		c.inWord = false
		return
	}

	c.lineRunes++

	if isWhitespace {
		c.inWord = false
		return
	}

	if !c.inWord {
		c.Words++
		c.inWord = true
	}
}

// ConsumeError treats an undecodable byte as a replacement rune.
func (c *Counters) ConsumeError(bytes int) {
	if bytes <= 0 {
		bytes = 1
	}
	c.Bytes += bytes
	c.Runes++
	c.lineRunes++
	c.lineDirty = true
	if !c.inWord {
		c.Words++
		c.inWord = true
	}
}

// Finish should be called after EOF to flush line statistics.
func (c *Counters) Finish() {
	c.updateMaxLine()
}

func (c *Counters) updateMaxLine() {
	if !c.lineDirty {
		return
	}
	if c.lineRunes > c.MaxLineLen {
		c.MaxLineLen = c.lineRunes
	}
	c.lineDirty = false
}

// Snapshot returns a copy of the public fields.
func (c *Counters) Snapshot() Snapshot {
	return Snapshot{
		Bytes:      c.Bytes,
		Runes:      c.Runes,
		Lines:      c.Lines,
		Words:      c.Words,
		MaxLineLen: c.MaxLineLen,
	}
}

// Snapshot captures final counter values.
type Snapshot struct {
	Bytes      int
	Runes      int
	Lines      int
	Words      int
	MaxLineLen int
}

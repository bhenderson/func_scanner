// Package func_scanner is modeled after text/scanner. It takes a simple []byte
// instead of an io.Reader. It too is for UTF-8-encoded text. A user defines
// their own split function which groups successive runes together. Any invalid
// UTF-8 sequence is still given to the split function.
//
// Basic Usage:
//
//	input := []byte("hi 123")
//
//	// this will split +input+ into "h", "i", " ", "123"
//	split := func(ch rune) rune {
//		// usually you'll group like runes with negative numbers and just
//		// return ch for anything else as a default
//		if ch >= '0' || ch <= '9' { // digit
//			return -1
//		}
//		return ch
//	}
//
//	s := Init(input, split)
//
//	for s.Scan() {
//		// s.Tok() returns whatever SplitFunc returned.
//		// s.Text() returns the string of grouped runes.
//	}
package func_scanner

import (
	"unicode/utf8"
)

// SplitFunc is a type for a function that takes a rune as input and returns a
// rune. Successive runes are grouped together. A return value of '0' causes
// scanning to stop.
type SplitFunc func(rune) rune

// Scanner is a struct type for the scanner to split a byte buffer into
// successive "similar" runes according to SplitFunc.
type Scanner struct {
	tok   rune
	split SplitFunc
	buf   []byte
	len   int // cache len(buf). TODO Although, it could change.
	i     int
	end   int

	ntok rune
	size int
}

// Init initializes a new Scanner. SplitFunc can be set here or at any time. It
// can also be changed inbetween scans.
func Init(p []byte, split ...SplitFunc) (s *Scanner) {
	s = &Scanner{
		buf: p,
		len: len(p),
	}

	if len(split) > 0 {
		s.Split(split[0])
	}

	return
}

// Split sets the SplitFunc.
func (s *Scanner) Split(split SplitFunc) {
	s.split = split
}

// Scan returns true until the whole []byte has been scanned or SplitFunc
// returns 0
func (s *Scanner) Scan() bool {
	if s.tok == 0 { //first one
		s.next()
	}

	s.i = s.end

	if s.i >= s.len { // EOF
		return false
	}

	s.next()

	for s.tok == s.ntok && s.end < s.len {
		s.next()
	}

	// 0 is treated as an invalid return value as it's the nil rune
	if s.tok == 0 {
		return false
	}

	return true
}

// Text returns the current group of similar runes as a string.
func (s *Scanner) Text() string {
	return string(s.buf[s.i:s.end])
}

// Tok is whatever SplitFunc returns.
func (s *Scanner) Tok() rune {
	return s.tok
}

func (s *Scanner) next() {
	s.end += s.size
	s.tok = s.ntok

	ch, size := utf8.DecodeRune(s.buf[s.end:])
	s.size = size
	s.ntok = s.split(ch)
}

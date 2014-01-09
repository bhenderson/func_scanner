package func_scanner

import (
	"unicode/utf8"
)

type SplitFunc func(rune) rune

type Scanner struct {
	tok   rune
	split SplitFunc
	buf   []byte
	len   int
	i     int
	end   int

	ntok rune
	size int
}

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

func (s *Scanner) Split(split SplitFunc) {
	s.split = split
}

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

func (s *Scanner) Text() string {
	return string(s.buf[s.i:s.end])
}

// If tok is whatever SplitFunc returns. If 0, then scanning stops.
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

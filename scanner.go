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

func Init(p []byte, split SplitFunc) (s *Scanner) {
	s = &Scanner{
		buf:   p,
		len:   len(p),
		split: split,
	}
	s.next() // initialize the first one.

	return
}

func (s *Scanner) Scan() bool {
	s.i = s.end

	if s.i >= s.len { // EOF
		return false
	}

	s.next()

	for s.tok == s.ntok {
		s.next()
	}

	return true
}

func (s *Scanner) Text() string {
	return string(s.buf[s.i:s.end])
}

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

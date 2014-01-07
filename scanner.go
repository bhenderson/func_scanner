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
	start int
	end   int

	ntok   rune
	nstart int
	nend   int
}

func Init(p []byte, split SplitFunc) (s *Scanner) {
	s = &Scanner{
		buf:   p,
		len:   len(p),
		split: split,
	}

	return
}

func (s *Scanner) Scan() bool {
	s.i = s.nstart

	if s.i >= s.len { // EOF
		return false
	}

	s.replace()

	if s.nstart == 0 { // first one
		s.next()
		s.replace()
	}

	s.next()

	for s.tok == s.ntok {
		s.replace()
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
	ch, size := utf8.DecodeRune(s.buf[s.nstart:])
	s.nend = s.nstart + size
	s.ntok = s.split(ch)
}

func (s *Scanner) replace() {
	s.start, s.end, s.tok, s.nstart = s.nstart, s.nend, s.ntok, s.nend
}

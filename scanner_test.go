package func_scanner

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"unicode"
)

func TestScanner(t *testing.T) {
	input := []byte("abc.123-a1b1")
	exp := []interface{}{"abc", ".", 123, "-", "a", 1, "b", 1}
	var act []interface{}

	split := func(ch rune) (t rune) {
		switch {
		case ch == '.' || ch == '-':
			t = ch
		case unicode.IsDigit(ch):
			t = -1
		case unicode.IsLetter(ch):
			t = -2
		}

		return
	}

	scanner := Init(input, split)

	for scanner.Scan() {
		switch scanner.Tok() {
		default: //letters and everything else
			act = append(act, scanner.Text())
		case -1: //digit
			d, err := strconv.Atoi(scanner.Text())

			if err != nil {
				t.Fatal(err)
			}

			act = append(act, d)
		}
	}

	assert.Equal(t, exp, act)
}

func TestInvalidSplit(t *testing.T) {
	input := []byte("hello")

	split := func(ch rune) rune {
		return 0
	}

	s := Init(input, split)

	assert.False(t, s.Scan())
}

func TestSplitReturnsZero(t *testing.T) {
	input := []byte("hi world")
	exp := []string{"hi"}
	var act []string

	split := func(ch rune) rune {
		if ch == ' ' { // kill scanning
			return 0
		}

		return -1
	}

	s := Init(input, split)

	for s.Scan() {
		act = append(act, s.Text())
	}

	assert.Equal(t, exp, act)
}

func TestMultiByte(t *testing.T) {
	input := []byte("hi 世") // don't know what this says, sorry
	exp := []string{"hi", " ", "世"}
	var act []string

	split := func(ch rune) (t rune) {
		if unicode.IsLetter(ch) {
			t = -1
		} else {
			t = -2
		}
		return
	}

	s := Init(input, split)

	for s.Scan() {
		act = append(act, s.Text())
	}

	assert.Equal(t, exp, act)
}

func TestInvalidBytes(t *testing.T) {
	input := []byte{104, 101, 228, 184, 108, 108, 111} // he<invalid>llo
	exp := []string{"he", "\xe4\xb8", "llo"}
	var act []string

	split := func(ch rune) (t rune) {
		if unicode.IsLetter(ch) {
			t = -1
		} else {
			t = -2 // error
		}
		return
	}

	s := Init(input, split)

	for s.Scan() {
		act = append(act, s.Text())
	}

	assert.Equal(t, exp, act)
}

func TestScannerDebug(t *testing.T) {
	input := []byte("abc.123-a1b1")

	split := func(ch rune) (t rune) {
		switch {
		case ch == '.' || ch == '-':
			t = ch
		case unicode.IsDigit(ch):
			t = -1
		case unicode.IsLetter(ch):
			t = -2
		}

		return
	}

	scanner := Init(input, split)

	for i := 0; i < 13; i++ {
		t.Log(scanner.Scan(), scanner)
	}
}

func TestSetSplit(t *testing.T) {
	input := []byte("hi")

	split := func(ch rune) rune {
		return -1
	}

	s := Init(input)
	s.Split(split)

	assert.Equal(t, -1, s.split('h'))
}

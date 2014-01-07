package func_scanner

import (
	"reflect"
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

	assertEqual(t, exp, act)
}

func TestInvalidBytes(t *testing.T) {
	input := []byte{104, 101, 228, 184, 108, 108, 111} // he<invalid>llo

	split := func(ch rune) (t rune) {
		if ch == '世' {
			t = -1
		}
		return
	}

	s := Init(input, split)

	refute(t, s.Scan())
}

func TestScannerTest(t *testing.T) {
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

func assertEqual(t *testing.T, expected interface{}, actual interface{}, msg ...string) {
	assert(t, reflect.DeepEqual(expected, actual),
		"%v\nexpected\n\t(%T)%#v\nto be equal to\n\t(%T)%#v",
		msg, expected, expected, actual, actual)
}

func assert(t *testing.T, act bool, msg ...interface{}) {
	if !act {
		switch v := msg[0].(type) {
		default:
			t.Fatal(msg)
		case string:
			t.Fatalf(v, msg[1:]...)
		}
	}
}

func refute(t *testing.T, act bool, msg ...interface{}) {
	assert(t, !act, msg...)
}

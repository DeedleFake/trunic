package trunic

import (
	"fmt"
	"iter"
	"strings"
	"unicode/utf8"
)

// Normalize returns a copy of text with all unsupported characters
// removed.
func Normalize(text string) string {
	var buf strings.Builder
	for prefix := range validPrefixes(text) {
		buf.WriteString(prefix)
	}
	return buf.String()
}

func validPrefixes(text string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for len(text) > 0 {
			prefix, after, ok := cutValidPrefix(text)
			if !ok {
				_, i := utf8.DecodeRuneInString(text)
				text = text[i:]
				continue
			}

			if !yield(prefix) {
				return
			}
			text = after
		}
	}
}

func cutValidPrefix(text string) (prefix, after string, ok bool) {
	after, ok = strings.CutPrefix(text, " ")
	if ok {
		return " ", after, true
	}

	for _, p := range prefixes {
		after, ok = strings.CutPrefix(text, p)
		if ok {
			return p, after, true
		}
	}

	return "", text, false
}

// Runes returns a sequence of rune sections suitable for passing to
// [Renderer.AppendRune].
//
// The yielded slice should not be retained from one iteration to the
// next.
func Runes(text string) iter.Seq[[]string] {
	return func(yield func([]string) bool) {
		var y []string
		for prefix := range validPrefixes(text) {
			switch len(y) {
			case 0:
				y = append(y, prefix)
				if !IsLetter(prefix) {
					if prefix == " " {
						y = y[:0]
					}
					if !yield(y) {
						return
					}
					y = y[:0]
				}

			case 1:
				if !IsLetter(prefix) {
					if !yield(y) {
						return
					}

					y[0] = prefix
					if prefix == " " {
						y = y[:0]
					}
					if !yield(y) {
						return
					}
					y = y[:0]
					continue
				}

				switch {
				case IsVowel(y[0]) && IsConsonant(prefix):
					y = append(y, prefix, "*")

				case IsConsonant(y[0]) && IsVowel(prefix):
					y = append(y, prefix)
				}

				if !yield(y) {
					return
				}
				if len(y) > 1 {
					y = y[:0]
					continue
				}
				y[0] = prefix

			default:
				panic(fmt.Errorf("invalid len(y): %v", len(y)))
			}
		}

		if len(y) != 0 {
			yield(y)
		}
	}
}

// IsConsonant returns true if ph is a consonant.
func IsConsonant(ph string) bool {
	_, ok := consonants[ph]
	return ok
}

// IsVowel returns true if ph is a vowel.
func IsVowel(ph string) bool {
	_, ok := vowels[ph]
	return ok
}

// IsLetter returns true if ph is a consonant or a vowel.
func IsLetter(ph string) bool {
	return IsConsonant(ph) || IsVowel(ph)
}

// IsSymbol returns true if ph is a symbol, such as punctuation or a
// reversing circle.
func IsSymbol(ph string) bool {
	_, ok := symbols[ph]
	return ok
}

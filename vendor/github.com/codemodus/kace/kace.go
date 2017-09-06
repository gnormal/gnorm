// Package kace provides common case conversion functions which take into
// consideration common initialisms.
package kace

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/codemodus/kace/ktrie"
)

const (
	kebabDelim = '-'
	snakeDelim = '_'
)

var (
	ciTrie *ktrie.KTrie
)

func init() {
	var err error
	if ciTrie, err = ktrie.NewKTrie(ciMap); err != nil {
		panic(err)
	}
}

// Camel returns a camelCased string.
func Camel(s string) string {
	return camelCase(ciTrie, s, false)
}

// Pascal returns a PascalCased string.
func Pascal(s string) string {
	return camelCase(ciTrie, s, true)
}

// Kebab returns a kebab-cased string with all lowercase letters.
func Kebab(s string) string {
	return delimitedCase(s, kebabDelim, false)
}

// KebabUpper returns a KEBAB-CASED string with all upper case letters.
func KebabUpper(s string) string {
	return delimitedCase(s, kebabDelim, true)
}

// Snake returns a snake_cased string with all lowercase letters.
func Snake(s string) string {
	return delimitedCase(s, snakeDelim, false)
}

// SnakeUpper returns a SNAKE_CASED string with all upper case letters.
func SnakeUpper(s string) string {
	return delimitedCase(s, snakeDelim, true)
}

// Kace provides common case conversion methods which take into
// consideration common initialisms set by the user.
type Kace struct {
	t *ktrie.KTrie
}

// New returns a pointer to an instance of kace loaded with a common
// initialsms trie based on the provided map. Before conversion to a
// trie, the provided map keys are all upper cased.
func New(initialisms map[string]bool) (*Kace, error) {
	ci := initialisms
	if ci == nil {
		ci = map[string]bool{}
	}

	ci = sanitizeCI(ci)

	t, err := ktrie.NewKTrie(ci)
	if err != nil {
		return nil, fmt.Errorf("kace: cannot create trie: %s", err)
	}

	k := &Kace{
		t: t,
	}

	return k, nil
}

// Camel returns a camelCased string.
func (k *Kace) Camel(s string) string {
	return camelCase(k.t, s, false)
}

// Pascal returns a PascalCased string.
func (k *Kace) Pascal(s string) string {
	return camelCase(k.t, s, true)
}

// Snake returns a snake_cased string with all lowercase letters.
func (k *Kace) Snake(s string) string {
	return delimitedCase(s, snakeDelim, false)
}

// SnakeUpper returns a SNAKE_CASED string with all upper case letters.
func (k *Kace) SnakeUpper(s string) string {
	return delimitedCase(s, snakeDelim, true)
}

// Kebab returns a kebab-cased string with all lowercase letters.
func (k *Kace) Kebab(s string) string {
	return delimitedCase(s, kebabDelim, false)
}

// KebabUpper returns a KEBAB-CASED string with all upper case letters.
func (k *Kace) KebabUpper(s string) string {
	return delimitedCase(s, kebabDelim, true)
}

func camelCase(t *ktrie.KTrie, s string, ucFirst bool) string {
	rs := []rune(s)
	d := 0
	prev := rune(-1)

	for i := 0; i < len(rs); i++ {
		r := rs[i]

		if unicode.IsLetter(r) {
			isToUpper := isToUpperInCamel(prev, r, ucFirst)

			tprev, skip := updateRunes(rs, i, d, t, isToUpper)
			if skip > 0 {
				i += skip
				prev = tprev
				continue
			}

			prev = updateRune(rs, i, d, isToUpper)
			continue
		}

		if unicode.IsDigit(r) {
			prev = updateRune(rs, i, d, false)
			continue
		}

		prev = r
		d++
	}

	return string(rs[:len(rs)-d])
}

func updateRune(rs []rune, i, delta int, upper bool) rune {
	r := rs[i]

	targ := i - delta
	if targ < 0 || i > len(rs)-1 {
		panic("this function has been used or designed incorrectly")
	}

	fn := unicode.ToLower
	if upper {
		fn = unicode.ToUpper
	}

	rs[targ] = fn(r)

	return r
}

func updateRunes(rs []rune, i, delta int, t *ktrie.KTrie, upper bool) (rune, int) {
	r := rs[i]
	ct := 0

	for j := t.MaxDepth(); j >= t.MinDepth(); j-- {
		if i+j <= len(rs) && t.FindAsUpper(rs[i:i+j]) {
			r = rs[i+j-1]
			ct = j - 1
			break
		}
	}

	if ct > 0 {
		for j := i; j <= i+ct; j++ {
			targ := j - delta
			if targ < 0 {
				panic("this function has been used or designed incorrectly")
			}

			fn := unicode.ToLower
			if upper {
				fn = unicode.ToUpper
			}

			rs[targ] = fn(rs[j])
		}
	}

	return r, ct
}

func isToUpperInCamel(prev, curr rune, ucFirst bool) bool {
	if prev == -1 {
		return ucFirst
	}

	if !unicode.IsLetter(prev) || unicode.IsUpper(curr) && unicode.IsLower(prev) {
		return true
	}

	return false
}

func delimitedCase(s string, delim rune, upper bool) string {
	buf := make([]rune, 0, len(s)*2)

	for i := len(s); i > 0; i-- {
		switch {
		case unicode.IsLetter(rune(s[i-1])):
			if i < len(s) && unicode.IsUpper(rune(s[i])) {
				if i > 1 && unicode.IsLower(rune(s[i-1])) || i < len(s)-2 && unicode.IsLower(rune(s[i+1])) {
					buf = append(buf, delim)
				}
			}

			buf = appendCased(buf, upper, rune(s[i-1]))

		case unicode.IsDigit(rune(s[i-1])):
			if i == len(s) || i == 1 || unicode.IsDigit(rune(s[i])) {
				buf = append(buf, rune(s[i-1]))
				continue
			}

			buf = append(buf, delim, rune(s[i-1]))

		default:
			if i == len(s) {
				continue
			}

			buf = append(buf, delim)
		}
	}

	reverse(buf)

	return string(buf)
}

func appendCased(rs []rune, upper bool, r rune) []rune {
	if upper {
		rs = append(rs, unicode.ToUpper(r))
		return rs
	}

	rs = append(rs, unicode.ToLower(r))

	return rs
}

func reverse(s []rune) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

var (
	// github.com/golang/lint/blob/master/lint.go
	ciMap = map[string]bool{
		"ACL":   true,
		"API":   true,
		"ASCII": true,
		"CPU":   true,
		"CSS":   true,
		"DNS":   true,
		"EOF":   true,
		"GUID":  true,
		"HTML":  true,
		"HTTP":  true,
		"HTTPS": true,
		"ID":    true,
		"IP":    true,
		"JSON":  true,
		"LHS":   true,
		"QPS":   true,
		"RAM":   true,
		"RHS":   true,
		"RPC":   true,
		"SLA":   true,
		"SMTP":  true,
		"SQL":   true,
		"SSH":   true,
		"TCP":   true,
		"TLS":   true,
		"TTL":   true,
		"UDP":   true,
		"UI":    true,
		"UID":   true,
		"UUID":  true,
		"URI":   true,
		"URL":   true,
		"UTF8":  true,
		"VM":    true,
		"XML":   true,
		"XMPP":  true,
		"XSRF":  true,
		"XSS":   true,
	}
)

func sanitizeCI(m map[string]bool) map[string]bool {
	r := map[string]bool{}

	for k := range m {
		fn := func(r rune) rune {
			if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
				return -1
			}
			return r
		}

		k = strings.Map(fn, k)
		k = strings.ToUpper(k)

		if k == "" {
			continue
		}

		r[k] = true
	}

	return r
}

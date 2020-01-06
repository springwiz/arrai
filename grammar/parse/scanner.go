package parse

import (
	"fmt"
	"regexp"
	"strings"
)

type Scanner struct {
	src    string
	slice  string
	offset int
}

func NewScanner(src string) *Scanner {
	return &Scanner{src: src, slice: src, offset: 0}
}

func NewScannerAt(src string, offset, size int) *Scanner {
	return &Scanner{src: src, slice: src[offset : offset+size], offset: offset}
}

func (r Scanner) String() string {
	return r.slice
}

func (r Scanner) Format(state fmt.State, c rune) {
	if c == 'q' {
		fmt.Fprintf(state, "%q", r.String())
	} else {
		state.Write([]byte(r.String()))
	}
}

func (r Scanner) Context() string {
	return fmt.Sprintf("%s\033[1;31m%s\033[0m%s",
		r.src[:r.offset],
		r.slice,
		r.src[r.offset+len(r.slice):],
	)
}

func (r Scanner) Offset() int {
	return r.offset
}

func (r Scanner) Slice(a, b int) Scanner {
	return Scanner{
		src:    r.src,
		slice:  r.slice[a:b],
		offset: r.offset + a,
	}
}

func (r Scanner) Skip(i int) Scanner {
	return r.Slice(i, len(r.slice))
}

func (r Scanner) Eat(i int) Scanner {
	r.slice = r.slice[:i]
	return r
}

func (r Scanner) EatString(s string) (Scanner, bool) {
	if strings.HasPrefix(r.String(), s) {
		return r.Eat(len(s)), true
	}
	return Scanner{}, false
}

// EatRegexp eats the text matchin a regexp, as much of the who match and any
// captured groups as will fit into captures. Returns n as the number of
// captures set and ok iff a match was found.
func (r Scanner) EatRegexp(re *regexp.Regexp, captures []Scanner) (n int, remaning Scanner, ok bool) {
	if loc := re.FindStringSubmatchIndex(r.String()); loc != nil {
		if loc[0] != 0 {
			panic(`re not \A-anchored`)
		}
		n = len(loc) / 2
		if len(captures) > n {
			captures = captures[:n]
		}
		for i := range captures {
			captures[i] = r.Slice(loc[2*i], loc[2*i+1])
		}
		return n, r.Skip(loc[1]), true
	}
	return 0, Scanner{}, false
}

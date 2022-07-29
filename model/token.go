package model

import "regexp"

type TokenDef struct {
	Type    string
	Pattern *regexp.Regexp
}

type Token struct {
	Type  string
	Value string
	Line  int
	Col   int
}

func (t *Token) isAfter(other *Token) bool {
	if t.Line < other.Line {
		return false
	} else if t.Line > other.Line {
		return true
	} else {
		return t.Col > other.Col
	}
}

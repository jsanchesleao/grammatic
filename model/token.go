package model

import "regexp"

// Holds the data necessary for the lexer to output tokens of the defined type
type TokenDef struct {
	Type    string
	Pattern *regexp.Regexp
}

// Holds a chunk of the original parsed input, as well as the matched token type and position
type Token struct {
	Type  string
	Value string
	Line  int
	Col   int
}

// Checks if other token comes after the given token in the original input
func (t Token) IsAfter(other Token) bool {
	if t.Line < other.Line {
		return false
	} else if t.Line > other.Line {
		return true
	} else {
		return t.Col > other.Col
	}
}

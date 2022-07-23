package grammatic

import (
	"fmt"
	"regexp"
)

// TokenDef represents a definition to generate tokens of a particular type
// It requires a name, that will be present in the generated token for identification
// Also needs a regexp pattern, to check the input string.
// It's important that the regexp begins with a ^, to make it search the beginning of the string.
// There are some convenience regexp patterns exported from this package,
// like KeywordFormat and FloatTokenFormat.
type TokenDef struct {
	Name    string
	Pattern *regexp.Regexp
}

// Token represents a meaningful part of the input string, to be used later by the parser.
// The token carries the name defined in the matching TokenDef, as well as the string part
// that matches the pattern, and also line and column numbers, for better error recognition
// in the input string.
type Token struct {
	Name  string
	Value string
	Line  int
	Col   int
}

const TYPE_EOF = "TOKEN_EOF"

const DigitsTokenFormat = "^\\d+"
const IntTokenFormat = "^[123456789]\\d*"
const FloatTokenFormat = "^[123456789]\\d*\\.\\d+"
const KeywordFormat = "^(?i)[abcdefghijklmnopqrstuvwxyz][-_\\w]*"
const DoubleQuotedStringFormat = "^\"(\\\"|[^\\\"])*\""
const EmptySpaceFormat = "^\\s+"
const OperandFormat = "^[-+/*=]"
const OpenBracesFormat = "^(\\(|\\[|\\{)"
const CloseBracesFormat = "^(\\)|\\]|\\})"
const PunctuationFormat = "^[,;:.]"

// Convenience function to create a TokenDef, that compiles the regexp pattern.
// This is specially convenient when used with the provided patters.
// i.e grammatic.NewTokenDef("EmptySpace", grammatic.EmptySpaceFormat)
func NewTokenDef(name, pattern string) TokenDef {
	regex := regexp.MustCompile(pattern)
	return TokenDef{Name: name, Pattern: regex}
}

// Accepts an input string and a slice of TokenDef, and returns a slice of tokens when successful,
// or an error if anything goes wrong
// All characters existing in the input string must be recognizable in at least one of the
// tokendefs, otherwise an "Illegal Character" error will be generated.
// TokenDefs are tested in the order they are passed in to the function, and the first one
// that matches will generate a Token, so it's a good idea to place more specific TokenDefs
// first in the slice.
// At the end, if successful, this function will append an TYPE_EOF token, as the last one,
// indicating the end of the input. This will always be generated if no error was found.
func ExtractTokens(text string, tokendefs []TokenDef) ([]Token, error) {
	tokens := []Token{}
	line := 1
	col := 0
	index := 0
	var err error

	skips := 0

	for {
		if index >= len(text) {
			tokens = append(tokens, Token{Name: TYPE_EOF, Value: "", Line: line + 1, Col: 0})
			break
		}

		nextToken := Token{}
		if text[index] == '\n' {
			nextToken.Col = col + 1
			nextToken.Line = line
			col = 0
			line++
		} else {
			col++
			nextToken.Col = col
			nextToken.Line = line
		}

		if skips > 0 {
			skips--
			index++
			continue
		}

		remainingText := text[index:]
		hasToken := false
		for _, def := range tokendefs {
			if match := def.Pattern.FindString(remainingText); match != "" {
				nextToken.Name = def.Name
				nextToken.Value = match
				hasToken = true
				skips = len(match) - 1
				break
			}
		}
		if !hasToken {
			err = fmt.Errorf("Illegal character %q at line %d, column %d", string(text[index]), line, col)
			break
		} else {
			tokens = append(tokens, nextToken)
		}
		index++
	}

	return tokens, err
}

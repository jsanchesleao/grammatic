package grammatic

import (
	"fmt"
	"regexp"
)

type TokenDef struct {
	Name    string
	Pattern *regexp.Regexp
}

type Token struct {
	Name  string
	Value string
	Line  int
	Col   int
}

const DigitsTokenFormat = "^\\d+"
const IntTokenFormat = "^[123456789]\\d*"
const KeywordFormat = "^(?i)\\w[-_\\w]*"
const DoubleQuotedStringFormat = "^\"(\\\"|[^\\\"])*\""
const EmptySpaceFormat = "^\\s+"

func NewTokenDef(name, pattern string) TokenDef {
	regex := regexp.MustCompile(pattern)
	return TokenDef{Name: name, Pattern: regex}
}

func ExtractTokens(text string, tokendefs []TokenDef) ([]Token, error) {
	tokens := []Token{}
	line := 1
	col := 0
	index := 0
	var err error

	skips := 0

	for {
		if index >= len(text) {
			break
		}

		if text[index] == '\n' {
			col = 1
			line++
		} else {
			col++
		}

		if skips > 0 {
			skips--
			index++
			continue
		}

		remainingText := text[index:]
		nextToken := Token{}
		hasToken := false
		for _, def := range tokendefs {
			if match := def.Pattern.FindString(remainingText); match != "" {
				nextToken.Name = def.Name
				nextToken.Value = match
				nextToken.Line = line
				nextToken.Col = col
				hasToken = true
				skips = len(match) - 1
				break
			}
		}
		if !hasToken {
			err = fmt.Errorf("Illegal character at %d:%d", line, col)
			break
		} else {
			tokens = append(tokens, nextToken)
		}
		index++
	}

	return tokens, err

}

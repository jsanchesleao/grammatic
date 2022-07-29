package lexer

import (
	"fmt"
	"grammatic/model"
	"regexp"
)

const TYPE_EOF = "TOKEN_EOF"

const DigitsTokenFormat = "^\\d+"
const IntTokenFormat = "^[123456789]\\d*"
const FloatTokenFormat = "^[123456789]\\d*\\.\\d+"
const NumberTokenFormat = "^[123456789]\\d*(\\.\\d+)?"
const KeywordFormat = "^(?i)[abcdefghijklmnopqrstuvwxyz][-_\\w]*"
const DoubleQuotedStringFormat = "^\"(\\\"|[^\\\"])*?\""
const EmptySpaceFormat = "^\\s+"
const OperandFormat = "^[-+/*=]"
const OpenBracesFormat = "^(\\(|\\[|\\{)"
const CloseBracesFormat = "^(\\)|\\]|\\})"
const PunctuationFormat = "^[,;:.]"

func NewTokenDef(tokenType, pattern string) model.TokenDef {
	regex := regexp.MustCompile(pattern)
	return model.TokenDef{Type: tokenType, Pattern: regex}
}

func ExtractTokens(text string, tokendefs []model.TokenDef) ([]model.Token, error) {
	tokens := []model.Token{}
	line := 1
	col := 0
	index := 0
	var err error

	skips := 0

	for {
		if index >= len(text) {
			tokens = append(tokens, model.Token{Type: TYPE_EOF, Value: "", Line: line + 1, Col: 0})
			break
		}

		nextToken := model.Token{}
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
				nextToken.Type = def.Type
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

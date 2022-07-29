package model

import (
	"fmt"
)

type Rule struct {
	Type  string
	Check func([]Token) RuleResultIterator
}

type RuleResultIterator interface {
	Next() *RuleResult
	Done()
}

type FinalResultCandidateType struct {
	next   func() *RuleResult
	finish func()
}

func FinalResultCandidate(result RuleResult) FinalResultCandidateType {
	done := false
	return FinalResultCandidateType{
		next: func() *RuleResult {
			if done {
				return nil
			}
			return &result
		},
		finish: func() {
			done = true
		},
	}
}

func (f FinalResultCandidateType) Next() *RuleResult {
	result := f.next()
	f.finish()
	return result
}

func (f FinalResultCandidateType) Done() {
	f.finish()
}

type RuleResult struct {
	Match           *Node
	RemainingTokens []Token
	Error           *RuleError
}

type RuleError struct {
	RuleType string
	Token    Token
}

func (e *RuleError) GetError() error {
	return fmt.Errorf("Unexpected token %q at line %d, column %d", e.Token.Value, e.Token.Line, e.Token.Col)
}

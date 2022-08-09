package model

import (
	"fmt"
)

// Represents a Generator Rule, which has a type name and a verifying function
type Rule struct {
	Type  string
	Check func([]Token) RuleResultIterator
}

// Returned by a Rule, this will output RuleResults with the Next() method and nil after it's finished, or after Done() is called
type RuleResultIterator interface {
	Next() *RuleResult
	Done()
}

type FinalResultCandidateType struct {
	next   func() *RuleResult
	finish func()
}

// Represents a RuleResultIterator that emits one single result and then ends (gives only nil afterwards)
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

// Represents an intermediary state for the parsing process, with the last produced result or error, plus the remaining tokens to be parsed
type RuleResult struct {
	Match           *Node
	RemainingTokens []Token
	Error           *RuleError
}

// Represents an "Unexpected Token" error
type RuleError struct {
	RuleType string
	Token    Token
}

// Converts from a RuleError to a golang standard error type
func (e *RuleError) GetError() error {
	return fmt.Errorf("Unexpected token %q at line %d, column %d", e.Token.Value, e.Token.Line, e.Token.Col)
}

package parser

import (
	"github.com/jsanchesleao/grammatic/model"
)

const (
	NULL = iota
	REQUEST_NEXT
	REQUEST_DONE
	CONTROL_CONTINUE
	CONTROL_BREAK
)

type ResultStream struct {
	feed    chan *model.RuleResult
	control chan int
	output  chan *model.RuleResult
	done    bool

	NodeMapper func(model.Node) model.Node
}

func NewResultStream() *ResultStream {
	stream := ResultStream{
		feed:       make(chan *model.RuleResult),
		control:    make(chan int),
		output:     make(chan *model.RuleResult),
		done:       false,
		NodeMapper: func(node model.Node) model.Node { return node },
	}

	return &stream
}

func (s ResultStream) Continue() bool {
	control := <-s.control
	return control == CONTROL_CONTINUE
}

func (s ResultStream) Send(result *model.RuleResult) {
	if s.done {
		return
	}
	if result != nil && result.Match != nil {
		match := s.NodeMapper(*result.Match)
		s.output <- &model.RuleResult{
			Match:           &match,
			RemainingTokens: result.RemainingTokens,
			Error:           result.Error,
		}
	} else {
		s.output <- result
	}
}

func (s *ResultStream) Next() *model.RuleResult {
	if s.done {
		return nil
	}

	defer func() {
		recover()
	}()

	s.control <- CONTROL_CONTINUE
	value := <-s.output
	return value
}

func (s *ResultStream) Done() {
	if !s.done {
		s.done = true
		close(s.output)
		close(s.control)
	}
}

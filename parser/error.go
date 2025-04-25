package parser

import "fmt"

type ParserError struct {
	Msg    string
	Detail string
	Hint   string
	Line   int // optional for now
	Column int // optional for now
}

func (e *ParserError) Error() string {
	msg := fmt.Sprintf("\n[Parser Error] %s", e.Msg)
	if e.Detail != "" {
		msg += fmt.Sprintf("\n  → %s", e.Detail)
	}
	if e.Hint != "" {
		msg += fmt.Sprintf("\n  💡 Hint: %s", e.Hint)
	}
	return msg
}

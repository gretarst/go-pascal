package interpreter

import "fmt"

type PascalError struct {
	Msg    string
	Detail string
	Hint   string
}

func (e *PascalError) Error() string {
	msg := fmt.Sprintf("\n[Pascal Error] %s", e.Msg)
	if e.Detail != "" {
		msg += fmt.Sprintf("\n  → %s", e.Detail)
	}
	if e.Hint != "" {
		msg += fmt.Sprintf("\n  💡 Hint: %s", e.Hint)
	}
	return msg
}

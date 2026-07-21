package http

import "fmt"

type Error struct {
	Code    int
	Reason  string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s", e.Code, e.Reason, e.Message)
}

func NewError(code int, reason, message string) *Error {
	return &Error{
		Code:    code,
		Reason:  reason,
		Message: message,
	}
}

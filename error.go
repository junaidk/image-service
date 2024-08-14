package imageapi

import (
	"errors"
	"fmt"
)

const (
	ERRCONFLICT = "conflict"
	ERRINTERNAL = "internal"
	ERRINVALID  = "invalid"
	ERRNOTFOUND = "not_found"
)

type Error struct {
	// Machine-readable error code.
	Code string

	// Human-readable error message.
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code=%s message=%s", e.Code, e.Message)
}

func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}
	return ERRINTERNAL
}

// ErrorMessage unwraps an application error and returns its message.
// Non-application errors always return "Internal error".
func ErrorMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}
	return "Internal error."
}

func Errorf(code string, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

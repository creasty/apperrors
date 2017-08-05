package apperrors

import (
	"errors"
	"fmt"
)

// Error is an error that has an application context
type Error struct {
	Err        error
	Message    string
	StatusCode int
	Report     bool
	StackTrace StackTrace
}

// New returns an error with the supplied message
func New(str string) error {
	return &Error{
		Err:        errors.New(str),
		StackTrace: newStackTrace(0),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error
func Errorf(format string, args ...interface{}) error {
	return &Error{
		Err:        fmt.Errorf(format, args...),
		StackTrace: newStackTrace(0),
	}
}

// Error implements error interface
func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return e.Err.Error()
}

// Copy creates a copy of the current object
func (e *Error) Copy() *Error {
	return &Error{
		Err:        e.Err,
		Message:    e.Message,
		StatusCode: e.StatusCode,
		Report:     e.Report,
		StackTrace: e.StackTrace,
	}
}

// Wrap returns an error annotated with a stack trace.
// If err is nil, Wrap returns nil.
func Wrap(err error) error {
	if err == nil {
		return nil
	}

	return wrap(err)
}

func wrap(err error) *Error {
	pkgErr := extractPkgError(err)

	if appErr, ok := pkgErr.Err.(*Error); ok {
		return appErr
	}

	stackTrace := pkgErr.StackTrace
	if stackTrace == nil {
		stackTrace = newStackTrace(1)
	}

	var msg string
	if pkgErr.Message != pkgErr.Err.Error() {
		msg = pkgErr.Message
	}

	return &Error{
		Err:        pkgErr.Err,
		StackTrace: stackTrace,
		Message:    msg,
	}
}

// Unwrap extracts underlying apperrors.Error from an error
func Unwrap(err error) *Error {
	if appErr, ok := err.(*Error); ok {
		return appErr
	}

	return nil
}

// WithMessage wraps err if necessary, and sets a message to its context
func WithMessage(err error, msg string) error {
	if err == nil {
		return nil
	}

	appErr := wrap(err).Copy()
	appErr.Message = msg
	return appErr
}

// WithStatusCode wraps err if necessary, and sets a status code to its context
func WithStatusCode(err error, code int) error {
	if err == nil {
		return nil
	}

	appErr := wrap(err).Copy()
	appErr.StatusCode = code
	return appErr
}

// WithReport wraps err if necessary, and marks as a reportable
func WithReport(err error) error {
	if err == nil {
		return nil
	}

	appErr := wrap(err).Copy()
	appErr.Report = true
	return appErr
}

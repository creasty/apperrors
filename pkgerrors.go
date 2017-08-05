package apperrors

import (
	"fmt"
	"strconv"

	pkgerrors "github.com/pkg/errors"
)

type pkgError struct {
	Err        error
	Message    string
	StackTrace StackTrace
}

func extractPkgError(err error) pkgError {
	type stackTracer interface {
		StackTrace() pkgerrors.StackTrace
	}

	type causer interface {
		Cause() error
	}

	var st pkgerrors.StackTrace
	rootErr := err

	for {
		if stackTrace, ok := rootErr.(stackTracer); ok {
			st = stackTrace.StackTrace()
		}

		if cause, ok := rootErr.(causer); ok {
			rootErr = cause.Cause()
		} else {
			break
		}
	}

	var frames []Frame
	if st != nil {
		for _, t := range st {
			file := fmt.Sprintf("%s", t)
			line, _ := strconv.ParseInt(fmt.Sprintf("%d", t), 10, 64)
			funcName := fmt.Sprintf("%n", t)

			frames = append(frames, Frame{
				Func: funcName,
				Line: line,
				File: file,
			})
		}
	}

	return pkgError{
		Err:        rootErr,
		Message:    err.Error(),
		StackTrace: frames,
	}
}

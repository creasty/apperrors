package apperrors

// Option attaches contextual metadata to errors.
type Option func(*Error)

// WithMessage wraps the error and annotates with the message.
// If err is nil, it returns nil
func WithMessage(msg string) Option {
	return func(err *Error) {
		err.Message = msg
	}
}

// WithStatusCode wraps the error and annotates with the status code.
// If err is nil, it returns nil
func WithStatusCode(code int) Option {
	return func(err *Error) {
		err.StatusCode = code
	}
}

// WithReport wraps the error and annotates with the reportability.
// If err is nil, it returns nil
func WithReport() Option {
	return func(err *Error) {
		err.Report = true
	}
}

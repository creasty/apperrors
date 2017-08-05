apperrors
=========

[![Build Status](https://travis-ci.org/creasty/apperrors.svg?branch=master)](https://travis-ci.org/creasty/apperrors)

Better error handling solution especially for application server.

`apperrors` provides contextual metadata to errors.

- Stack trace
- Additional information
- Status code (for a HTTP server)
- Reportability (for an integration with error reporting service)


Why
---

Since `error` type in Golang is just an interface of [`Error()`](https://golang.org/ref/spec#Errors) method, it doesn't have contextual information such as stack trace, you cannot be sure where the error occurred in the first place.  
And because of that, it's pretty hard to debug.

### How different from [pkg/errors](https://github.com/pkg/errors)

:memo: `creasty/apperrors` supports `pkg/errors`. It reuses `pkg/errors`'s stack trace data of the innermost (root) error, and converts into `apperrors`'s data type.

TBA



Create an error
---------------

```go
func New(str string) error
```

New returns an error that formats as the given text.  
It also annotates the error with a stack trace from the point it was called

```go
func Errorf(format string, args ...interface{}) error
```

Errorf formats according to a format specifier and returns the string
as a value that satisfies error.  
It also annotates the error with a stack trace from the point it was called

```go
func Wrap(err error) error
```

Wrap returns an error annotated with a stack trace from the point it was called.  
It returns nil if err is nil

### Example: Creating a new error

```go
ok := emailRegexp.MatchString("invalid#email.addr")
if !ok {
	return apperrors.New("invalid email address")
}
```

### Example: Creating from an existing error

```go
_, err := ioutil.ReadAll(r)
if err != nil {
	return apperrors.Wrap(err)
}
```


Annotate an error
-----------------

```go
func WithMessage(err error, msg string) error
```

WithMessage wraps the error and annotates with the message.  
If err is nil, it returns nil

```go
func WithStatusCode(err error, code int) error
```

WithStatusCode wraps the error and annotates with the status code.  
If err is nil, it returns nil

```go
func WithReport(err error) error
```

WithReport wraps the error and annotates with the reportability.  
If err is nil, it returns nil

### Example: Adding all contexts

```go
_, err := ioutil.ReadAll(r)
if err != nil {
	return apperrors.WithReport(
		apperrors.WithStatusCode(
			apperrors.WithMessage(err, "read failed"),
			http.StatusBadRequest
		)
	)
}
```


Extract context from an error
-----------------------------

```go
func Unwrap(err error) *Error
```

Unwrap extracts an underlying \*apperrors.Error from an error.  
If the given error isn't eligible for retriving context from,
it returns nil

```go
type Error struct {
	// Err is the original error (you might call it the root cause)
	Err error
	// Message is an annotated description of the error
	Message string
	// StatusCode is a status code that is desired to be used for a HTTP response
	StatusCode int
	// Report represents whether the error should be reported to administrators
	Report bool
	// StackTrace is a stack trace of the original error
	// from the point where it was created
	StackTrace StackTrace
}
```

### Example

Here's a minimum executable example describing how `apperrors` works.

```go
package main

import (
	"errors"
	"github.com/creasty/apperrors"
	"github.com/k0kubun/pp"
)

func errFunc0() error {
	return errors.New("this is the root cause")
}
func errFunc1() error {
	return apperrors.Wrap(errFunc0())
}
func errFunc2() error {
	return apperrors.WithMessage(errFunc1(), "fucked up!")
}
func errFunc3() error {
	return apperrors.WithReport(apperrors.WithStatusCode(errFunc2(), 500))
}

func main() {
	err := errFunc3()
	pp.Println(err)
}
```

```sh-session
$ go run main.go
&apperrors.Error{
  Err:        &errors.errorString{s: "this is the root cause"},
  Message:    "fucked up!",
  StatusCode: 500,
  Report:     true,
  StackTrace: apperrors.StackTrace{
    apperrors.Frame{Func: "errFunc1", File: "tmp/main.go", Line: 13},
    apperrors.Frame{Func: "errFunc2", File: "tmp/main.go", Line: 16},
    apperrors.Frame{Func: "errFunc3", File: "tmp/main.go", Line: 19},
    apperrors.Frame{Func: "main", File: "tmp/main.go", Line: 23},
    apperrors.Frame{Func: "main", File: "runtime/proc.go", Line: 194},
    apperrors.Frame{Func: "goexit", File: "runtime/asm_amd64.s", Line: 2198},
  },
}
```

### Example: Server-side error reporting with [gin-gonic/gin](https://github.com/gin-gonic/gin)

Prepare a simple middleware and modify to satisfy your needs:

```go
package middleware

const contextError = "AppError"

func SetError(c *gin.Context, err error) {
	c.Set(contextError, err)
}

func GetError(c *gin.Context) error {
	if v, exists := c.Get(contextError); exists {
		if err, ok := v.(error); ok {
			return err
		}
	}

	return nil
}

func ReportError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := GetError(c)
		if err == nil {
			return
		}

		if appErr := apperrors.Unwrap(err); appErr != nil {
			if appErr.Report {
				// Send to an external service
			}

			if appErr.Message != "" {
				// Expose a message in the header
				c.Header("X-App-Error", appErr.Message)
			}

			if appErr.StatusCode != 0 {
				// Set status code accordingly
				c.Status(appErr.StatusCode)
			}
		}
	}
}
```

And then you can use like as follows.

```go
r := gin.Default()
r.Use(middleware.ReportError()) // Use the middleware handler

r.GET("/test", func(c *gin.Context) {
	err := doSomethingReallyComplex()
	if err != nil {
		middleware.SetError(c, err) // Neither `c.AbortWithError` nor `c.Error`
		return
	}

	c.AbortWithStatus(200)
})

r.Run()
```

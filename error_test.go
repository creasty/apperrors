package apperrors

import (
	"errors"
	"testing"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New("message")
	assert.Equal(t, "message", err.Error())

	appErr := Unwrap(err)
	assert.Equal(t, err.Error(), appErr.Err.Error())
	assert.Equal(t, "", appErr.Message)
	assert.NotEmpty(t, appErr.StackTrace)
	assert.Equal(t, "TestNew", appErr.StackTrace[0].Func)
}

func TestErrorf(t *testing.T) {
	err := Errorf("message %d", 123)
	assert.Equal(t, "message 123", err.Error())

	appErr := Unwrap(err)
	assert.Equal(t, err.Error(), appErr.Err.Error())
	assert.Equal(t, "", appErr.Message)
	assert.NotEmpty(t, appErr.StackTrace)
	assert.Equal(t, "TestErrorf", appErr.StackTrace[0].Func)
}

func TestWithMessage(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Wrap(nil, WithMessage("message"))
		assert.Equal(t, nil, err)
	})

	t.Run("bare", func(t *testing.T) {
		err0 := errors.New("original")

		err1 := Wrap(err0, WithMessage("message"))
		assert.Equal(t, "message", err1.Error())

		appErr := Unwrap(err1)
		assert.Equal(t, err0, appErr.Err)
		assert.Equal(t, err1.Error(), appErr.Message)
	})

	t.Run("already wrapped", func(t *testing.T) {
		err0 := errors.New("original")

		err1 := &Error{
			Err:        err0,
			Message:    "message 1",
			StatusCode: 400,
		}
		err2 := Wrap(err1, WithMessage("message 2"))
		assert.Equal(t, "message 2", err2.Error())

		{
			appErr := Unwrap(err1)
			assert.Equal(t, err0, appErr.Err)
			assert.Equal(t, err1.Error(), appErr.Message)
			assert.Equal(t, 400, appErr.StatusCode)
		}

		{
			appErr := Unwrap(err2)
			assert.Equal(t, err0, appErr.Err)
			assert.Equal(t, err2.Error(), appErr.Message)
			assert.Equal(t, 400, appErr.StatusCode)
		}
	})
}

func TestWithStatusCode(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Wrap(nil, WithStatusCode(200))
		assert.Equal(t, nil, err)
	})

	t.Run("bare", func(t *testing.T) {
		err0 := errors.New("original")

		err1 := Wrap(err0, WithStatusCode(200))

		appErr := Unwrap(err1)
		assert.Equal(t, err0, appErr.Err)
		assert.Equal(t, "", appErr.Message)
	})

	t.Run("already wrapped", func(t *testing.T) {
		err0 := errors.New("original")

		err1 := &Error{
			Err:        err0,
			Message:    "message 1",
			StatusCode: 400,
		}
		err2 := Wrap(err1, WithStatusCode(500))

		{
			appErr := Unwrap(err1)
			assert.Equal(t, err0, appErr.Err)
			assert.Equal(t, err1.Error(), appErr.Message)
			assert.Equal(t, 400, appErr.StatusCode)
		}

		{
			appErr := Unwrap(err2)
			assert.Equal(t, err0, appErr.Err)
			assert.Equal(t, err1.Error(), appErr.Message)
			assert.Equal(t, 500, appErr.StatusCode)
		}
	})
}

func TestWithReport(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Wrap(nil, WithReport())
		assert.Equal(t, nil, err)
	})

	t.Run("bare", func(t *testing.T) {
		err0 := errors.New("original")

		err1 := Wrap(err0, WithReport())

		appErr := Unwrap(err1)
		assert.Equal(t, err0, appErr.Err)
		assert.Equal(t, "", appErr.Message)
	})

	t.Run("already wrapped", func(t *testing.T) {
		err0 := errors.New("original")

		err1 := Wrap(err0, WithReport())
		err2 := Wrap(err1, WithReport())

		{
			appErr := Unwrap(err1)
			assert.Equal(t, err0, appErr.Err)
			assert.Equal(t, true, appErr.Report)
		}

		{
			appErr := Unwrap(err2)
			assert.Equal(t, err0, appErr.Err)
			assert.Equal(t, true, appErr.Report)
		}
	})
}

func TestUnwrap(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		appErr := Unwrap(nil)
		assert.Nil(t, appErr)
	})
}

func TestWrap(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		appErr := Wrap(nil)
		assert.Nil(t, appErr)
	})

	t.Run("bare", func(t *testing.T) {
		err0 := errors.New("original")

		err1 := wrapOrigin(err0)
		assert.Equal(t, "original", err1.Error())

		appErr := Unwrap(err1)
		assert.Equal(t, err0, appErr.Err)
		assert.Equal(t, "", appErr.Message)
		assert.NotEmpty(t, appErr.StackTrace)
		assert.Equal(t, "wrapOrigin", appErr.StackTrace[0].Func)
	})

	t.Run("already wrapped", func(t *testing.T) {
		err0 := errors.New("original")

		err1 := wrapOrigin(err0)
		err2 := wrapOrigin(err1)
		assert.Equal(t, "original", err2.Error())

		appErr := Unwrap(err2)
		assert.Equal(t, err0, appErr.Err)
		assert.Equal(t, "", appErr.Message)
		assert.NotEmpty(t, appErr.StackTrace)
		assert.Equal(t, "wrapOrigin", appErr.StackTrace[0].Func)
	})

	t.Run("with pkg/errors", func(t *testing.T) {
		t.Run("pkg/errors.New", func(t *testing.T) {
			err0 := pkgErrorsNew("original")

			err1 := wrapOrigin(err0)
			assert.Equal(t, "original", err1.Error())

			appErr := Unwrap(err1)
			assert.Equal(t, err0, appErr.Err)
			assert.Equal(t, "", appErr.Message)
			assert.NotEmpty(t, appErr.StackTrace)
			assert.Equal(t, "pkgErrorsNew", appErr.StackTrace[0].Func)
		})

		t.Run("pkg/errors.Wrap", func(t *testing.T) {
			err0 := errors.New("original")
			err1 := pkgErrorsWrap(err0, "message")

			err2 := wrapOrigin(err1)
			assert.Equal(t, "message: original", err2.Error())

			appErr := Unwrap(err2)
			assert.Equal(t, err0, appErr.Err)
			assert.Equal(t, "message: original", appErr.Message)
			assert.NotEmpty(t, appErr.StackTrace)
			assert.Equal(t, "pkgErrorsWrap", appErr.StackTrace[0].Func)
		})
	})
}

func TestAll(t *testing.T) {
	{
		appErr := Unwrap(errFunc3())
		assert.Equal(t, "e2: e1: e0", appErr.Message)
		assert.Equal(t, 0, appErr.StatusCode)
		assert.Equal(t, false, appErr.Report)
		assert.NotEmpty(t, appErr.StackTrace)
		assert.Equal(t, "errFunc1", appErr.StackTrace[0].Func)
	}

	{
		appErr := Unwrap(errFunc4())
		assert.Equal(t, "e4", appErr.Message)
		assert.Equal(t, 500, appErr.StatusCode)
		assert.Equal(t, true, appErr.Report)
		assert.NotEmpty(t, appErr.StackTrace)
		assert.Equal(t, "errFunc1", appErr.StackTrace[0].Func)
	}
}

func wrapOrigin(err error) error {
	return Wrap(err)
}

func errFunc0() error {
	return errors.New("e0")
}
func errFunc1() error {
	return pkgerrors.Wrap(errFunc0(), "e1")
}
func errFunc2() error {
	return pkgerrors.Wrap(errFunc1(), "e2")
}
func errFunc3() error {
	return Wrap(errFunc2())
}
func errFunc4() error {
	return Wrap(errFunc3(), WithMessage("e4"), WithStatusCode(500), WithReport())
}

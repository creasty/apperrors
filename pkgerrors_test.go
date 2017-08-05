package apperrors

import (
	"errors"
	"testing"

	pkgerrors "github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestExtractPkgError(t *testing.T) {
	t.Run("pkg/errors.New", func(t *testing.T) {
		err := pkgErrorsNew("message")

		pkgErr := extractPkgError(err)
		assert.NotNil(t, pkgErr)
		assert.Equal(t, "message", pkgErr.Message)
		assert.Equal(t, err, pkgErr.Err)
		assert.NotEmpty(t, pkgErr.StackTrace)
		assert.Equal(t, "pkgErrorsNew", pkgErr.StackTrace[0].Func)
	})

	t.Run("pkg/errors.Wrap", func(t *testing.T) {
		err0 := errors.New("error")
		err1 := pkgErrorsWrap(err0, "message")

		pkgErr := extractPkgError(err1)
		assert.NotNil(t, pkgErr)
		assert.Equal(t, "message: error", pkgErr.Message)
		assert.Equal(t, err0, pkgErr.Err)
		assert.NotEmpty(t, pkgErr.StackTrace)
		assert.Equal(t, "pkgErrorsWrap", pkgErr.StackTrace[0].Func)
	})
}

func pkgErrorsNew(msg string) error {
	return pkgerrors.New(msg)
}

func pkgErrorsWrap(err error, msg string) error {
	return pkgerrors.Wrap(err, msg)
}

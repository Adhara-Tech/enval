package exerrors

import (
	"fmt"

	goerrors "github.com/go-errors/errors"
)

func Wrap(err error, kind EnvalErrorKind) error {
	return EnvalError{
		innerError: goerrors.Wrap(err, 3),
		kind:       kind,
	}
}

func New(msg string, kind EnvalErrorKind) error {
	return EnvalError{
		innerError: goerrors.Errorf(msg),
		kind:       kind,
	}
}

func PrintError(err error) string {
	envalError, ok := err.(EnvalError)
	if ok {
		return envalError.Error()
	} else {
		wrappedError := Wrap(err, InternalEnvalErrorKind)
		return wrappedError.Error()
	}
}

type EnvalErrorKind int

const (
	InternalEnvalErrorKind = iota
	ToolDefinitionNotFoundEnvalErrorKind
	FlavorDefinitionNotFoundEnvalErrorKind
	UnknownParserEnvalErrorKind
	UnsupportedInputRawVersionEnvalErrorKind
	FieldVersionKeyNotFoundEnvalErrorKind
	InvalidCustomSpecsDirEnvalErrorKind
)

func IsEnvalErrorWithKind(err error, kind EnvalErrorKind) bool {
	envalError, ok := err.(EnvalError)
	if !ok {
		return false
	}

	return envalError.kind == kind
}

type EnvalError struct {
	innerError *goerrors.Error
	kind       EnvalErrorKind
}

func (err EnvalError) Error() string {
	if err.kind == InternalEnvalErrorKind {
		return fmt.Sprintf("%s\n%s", err.innerError.Error(), err.innerError.ErrorStack())
	} else {
		return err.innerError.Error()
	}
}

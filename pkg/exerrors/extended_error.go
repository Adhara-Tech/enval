package exerrors

import (
	"fmt"

	goerrors "github.com/go-errors/errors"
)

func Wrap(err error) error {
	return goerrors.Wrap(err, 3)
}

func New(msg string) error {
	return goerrors.Errorf(msg)
}

func ErrorStack(err error) string {
	goError, ok := err.(*goerrors.Error)
	if ok {
		return fmt.Sprintf("%s\n%s", goError.Error(), goError.ErrorStack())
	} else {
		return err.Error()
	}
}

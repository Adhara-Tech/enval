package exerrors

import goerrors "github.com/go-errors/errors"

func Wrap(err error) error {
	return goerrors.Wrap(err, 3)
}

func New(msg string) error {
	return goerrors.Errorf(msg)
}

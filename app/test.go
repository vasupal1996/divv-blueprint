package app

import (
	"github.com/pkg/errors"
)

func TestError() error {
	return errorLast()
}

func errorStart() error {
	return errors.New("error start")
}

func errorMiddle() error {
	return errors.Wrap(errorStart(), "error middle")
}

func errorLast() error {
	return errors.Wrap(errorMiddle(), "error last")
}

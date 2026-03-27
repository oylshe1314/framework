package errors

import (
	"errors"
	"fmt"
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}

type stringError string

func (err stringError) Error() string {
	return string(err)
}

func New(text string) error {
	return stringError(text)
}

func Error(vals ...any) error {
	return stringError(fmt.Sprint(vals...))
}

func Errorf(format string, vals ...any) error {
	return stringError(fmt.Sprintf(format, vals...))
}

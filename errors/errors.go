package errors

import (
	"errors"
	"fmt"
)

type StringError string

func (err StringError) Error() string {
	return string(err)
}

const (
	ErrUnknownError StringError = "未知错误"
)

func Error(args ...any) error {
	if len(args) == 0 {
		return ErrUnknownError
	} else {
		return StringError(fmt.Sprint(args...))
	}
}

func Errorf(format string, args ...any) error {
	return StringError(fmt.Sprintf(format, args...))
}

const StatusUnknown int = -1
const StatusSuccessful int = 0

type StatusError interface {
	error

	Status() int
}

type statusError struct {
	status  int
	message string
}

func (this *statusError) Error() string {
	return this.message
}

func (this *statusError) Status() int {
	return this.status
}

func Status(status int, message string) error {
	return &statusError{status: status, message: message}
}

func Is(err, tag error) bool {
	var se1, se2 StatusError
	if errors.As(err, &se1) && errors.As(err, &se2) {
		return se1.Status() == se2.Status()
	}
	return errors.Is(err, tag)
}

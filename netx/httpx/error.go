package httpx

import "fmt"

type statusError struct {
	status  int
	message string
}

func (this *statusError) Error() string {
	return fmt.Sprint("Status: ", this.status, ", Message: ", this.message)
}

func StatusError(status int, message string) error {
	return &statusError{status, message}
}

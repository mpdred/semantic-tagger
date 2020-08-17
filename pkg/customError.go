package pkg

import "fmt"

type ErrorDetails interface {
	Error() string
	Type() error
}

type errDetails struct {
	errType error
	details interface{}
}

func (err *errDetails) Error() string {
	return fmt.Sprintf("%v: %v", err.errType, err.details)
}

func (err *errDetails) Type() error {
	return err.errType
}

func NewErrorDetails(err error, details ...interface{}) ErrorDetails {
	return &errDetails{
		errType: err,
		details: details,
	}
}

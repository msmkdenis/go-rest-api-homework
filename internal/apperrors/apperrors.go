package apperrors

import (
	"errors"
	"fmt"
)

var ErrTaskNotFound = errors.New("task not found")

type ValueError struct {
	caller string
	err    error
}

func NewValueError(caller string, err error) error {
	return &ValueError{
		caller: caller,
		err:    err,
	}
}

func (v *ValueError) Error() string {
	return fmt.Sprintf("%s error: %s", v.caller, v.err)
}

func (v *ValueError) Unwrap() error {
	return v.err
}

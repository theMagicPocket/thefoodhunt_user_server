package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNoDocuments = errors.New("no document(s) found")
)

func ErrNoRequiredQueryParam(param string) error {
	return fmt.Errorf("%s query param is required", param)
}

func ErrNoRequiredParam(param string) error {
	return fmt.Errorf("%s param is required", param)
}

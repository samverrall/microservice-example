package app

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrForbidden    = errors.New("forbidden")
)

func NewInvalidInputErr(err error) error {
	return fmt.Errorf("%w: %w", ErrInvalidInput, err)
}

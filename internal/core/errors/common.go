package core_errors

import "errors"

var (
	ErrNotFound = errors.New("not found")

	ErrInvalidArgument = errors.New("invalid argument")

	ErrInternal = errors.New("internal error")

	ErrConflict = errors.New("conflict")
)

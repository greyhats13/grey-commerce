// Path: grey-user/internal/app/errors.go

package app

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidRequest   = errors.New("invalid request")
	ErrFailedToParse    = errors.New("failed to parse request")
	ErrFailedToValidate = errors.New("failed to validate")
	ErrInternal         = errors.New("internal server error")
)

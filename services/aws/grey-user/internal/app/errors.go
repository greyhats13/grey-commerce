// Path: grey-user/internal/app/errors.go

package errors

import "errors"

var (
	ErrNotFound         = errors.New("resource not found")
	ErrInvalidRequest   = errors.New("invalid request")
	ErrFailedToParse    = errors.New("failed to parse request body")
	ErrFailedToValidate = errors.New("failed to validate")
	ErrInternal         = errors.New("internal server error")
)

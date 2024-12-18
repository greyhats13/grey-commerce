// Path: services/aws/grey-user/internal/app/errors.go

package app

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrInvalidRequest = errors.New("invalid request")
	ErrInternal       = errors.New("internal server error")
)

package constants

import "errors"

var (
	ErrNotFound   = func(resource string) error { return errors.New(resource + " not found") }
	ErrForbidden  = errors.New("forbidden: You do not have permission to perform this action")
	ErrBadRequest = errors.New("bad request")
)
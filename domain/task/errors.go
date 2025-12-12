package task

import (
	"errors"
)

var (
	ErrInvalidID     = errors.New("task: invalid id")
	ErrInvalidStatus = errors.New("task: invalid status")
	ErrInvalidTitle  = errors.New("task: invalid title")
	ErrNotFound      = errors.New("task: not found")
)

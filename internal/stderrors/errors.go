package stderrors

import "errors"

var (
	ErrInvalidURL        = errors.New("invalid url")
	ErrServiceNotReady   = errors.New("service not ready")
	ErrBrowserNotStarted = errors.New("browser adapter not started")
	ErrNavigationFailed  = errors.New("navigation failed")
	ErrTimeout           = errors.New("operation timed out")
)

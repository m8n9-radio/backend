package listener

import "errors"

var (
	ErrInvalidUserID    = errors.New("invalid user ID")
	ErrInvalidTrackID   = errors.New("invalid track ID")
	ErrListenerExists   = errors.New("listener already tracked")
	ErrListenerNotFound = errors.New("listener not found")
)

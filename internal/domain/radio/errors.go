package radio

import "errors"

var (
	ErrRadioUnavailable = errors.New("radio stream unavailable")
	ErrNoActiveStream   = errors.New("no active stream")
)

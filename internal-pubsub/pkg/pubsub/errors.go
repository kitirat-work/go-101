package pubsub

import "errors"

var (
	ErrClosed = errors.New("pubsub: bus is closed")
)

package exceptions

import "errors"

var (
	ErrEmailAlreadySubscribed = errors.New("email already subscribed")
)

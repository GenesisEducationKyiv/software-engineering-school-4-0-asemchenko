package exceptions

import "errors"

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrSubscriptionTimeout = errors.New("subscription timeout")
)

func ErrSubscriptionFailed(reason string) error {
	return errors.New("subscription failed: " + reason)
}

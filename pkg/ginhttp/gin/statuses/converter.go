package ginhttp

import (
	"errors"
)

// toError converts any error type to a standard error
func toError(err interface{}) error {
	switch e := err.(type) {
	case error:
		return e
	case string:
		return errors.New(e)
	default:
		return errors.New("unknown error")
	}
}

package exceptions

import "github.com/alexandria-oss/common-go/exception"

var (
	// ErrInvalidEvent an event did not comply with domain rules
	ErrInvalidEvent = exception.NewFieldFormat("event", "valid event")
	// ErrInvalidJSONEvent the given JSON event was not valid
	ErrInvalidJSONEvent = exception.NewFieldFormat("event", "valid json event")
)

package exceptions

import (
	"github.com/alexandria-oss/common-go/exception"
)

var (
	//	ErrBelowMinDuration is returned whenever a total duration is below 10 minutes
	ErrBelowMinDuration = exception.NewFieldRange("total_duration", "10 minutes", "48 hours")
	//	ErrAboveMaxDuration is returned whenever a total duration is above 48 hours
	ErrAboveMaxDuration = exception.NewFieldRange("total_duration", "10 minutes", "48 hours")
	//	ErrInvalidTotalDuration is returned whenever a total duration has an invalid value
	ErrInvalidTotalDuration = exception.NewFieldFormat("total_duration", "minutes (e.g. 10m)")
	//	ErrEmptyActivity is returned whenever an empty activity id is given
	ErrEmptyActivity = exception.NewRequiredField("activity_id")
	// ErrOccurrenceNotFound is returned whenever an occurrence was not found
	ErrOccurrenceNotFound = exception.NewNotFound("occurrence")
)

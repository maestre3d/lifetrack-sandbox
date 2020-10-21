package exceptions

import (
	"github.com/alexandria-oss/common-go/exception"
)

var (
	//	ErrBelowMinDuration total duration is below 10 minutes
	ErrBelowMinDuration = exception.NewFieldRange("total_duration", "10 minutes", "48 hours")
	//	ErrAboveMaxDuration total duration is above 48 hours
	ErrAboveMaxDuration = exception.NewFieldRange("total_duration", "10 minutes", "48 hours")
	//	ErrInvalidTotalDuration total duration has an invalid value
	ErrInvalidTotalDuration = exception.NewFieldFormat("total_duration", "minutes (e.g. 10m)")
	//	ErrEmptyActivityID empty activity id was given
	ErrEmptyActivityID = exception.NewRequiredField("activity_id")
	// ErrOccurrenceNotFound occurrence was not found
	ErrOccurrenceNotFound = exception.NewNotFound("occurrence")
	// ErrInvalidOccurrenceFilter occurrence fetch filter has no valid values
	ErrInvalidOccurrenceFilter = exception.NewFieldFormat("occurrence_filter", "occurrence__id, "+
		"activity_id, page_size, page_token")
)

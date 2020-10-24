package exceptions

import "github.com/alexandria-oss/common-go/exception"

var (
	// ErrEmptyCategoryID category field is empty
	ErrEmptyCategoryID = exception.NewRequiredField("category_id")
	// ErrBelowTitleLength title field is below 2 characters
	ErrBelowTitleLength = exception.NewFieldRange("title", "2", "512")
	// ErrActivityMarshaling problem occurred while activity marshaling
	ErrActivityMarshaling = exception.NewFieldFormat("activity", "json")
	// ErrActivityNotFound activity was not found
	ErrActivityNotFound = exception.NewNotFound("activity")
	// ErrInvalidActivityFilter activity fetch filter has no valid values
	ErrInvalidActivityFilter = exception.NewFieldFormat("activity_filter", "activity_id, "+
		"category_id, title, page_size, page_token")
)

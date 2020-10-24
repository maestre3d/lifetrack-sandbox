package exceptions

import "github.com/alexandria-oss/common-go/exception"

var (
	// ErrEmptyCategoryID category field is empty
	ErrEmptyCategoryID = exception.NewRequiredField("category")
	// ErrBelowTitleLength title field is below 2 characters
	ErrBelowTitleLength = exception.NewFieldRange("title", "2", "512")
	// ErrActivityMarshaling problem occurred while activity marshaling
	ErrActivityMarshaling = exception.NewFieldFormat("activity", "json")
)

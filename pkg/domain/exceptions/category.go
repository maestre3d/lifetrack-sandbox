package exceptions

import "github.com/alexandria-oss/common-go/exception"

var (
	// ErrEmptyUserID user field is empty
	ErrEmptyUserID = exception.NewRequiredField("user_id")
	// ErrBelowTitleLength name field is below 2 characters
	ErrBelowNameLength = exception.NewFieldRange("name", "2", "512")
	//	ErrInvalidTargetTime target time has an invalid value
	ErrInvalidTargetTime = exception.NewFieldFormat("target_time", "minutes (e.g. 10m)")
	// ErrCategoryMarshaling problem occurred while category marshaling
	ErrCategoryMarshaling = exception.NewFieldFormat("category", "json")
	// ErrCategoryNotFound category was not found
	ErrCategoryNotFound = exception.NewNotFound("category")
	// ErrInvalidCategoryFilter category fetch filter has no valid values
	ErrInvalidCategoryFilter = exception.NewFieldFormat("category_filter", "category_id, "+
		"user_id, title, page_size, page_token")
)

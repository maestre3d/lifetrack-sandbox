package value

import (
	"github.com/alexandria-oss/common-go/exception"
)

// Title text formatted using specific standards
type Title struct {
	value string

	fieldName string
}

// NewTitle creates a valid Title
func NewTitle(field, title string) (*Title, error) {
	t := new(Title)
	t.setFieldName(field)
	if err := t.Rename(title); err != nil {
		return nil, err
	}

	return t, nil
}

// NewTitleFromPrimitive creates a Title without validating for marshaling purposes
func NewTitleFromPrimitive(field, title string) *Title {
	t := &Title{
		value: title,
	}
	t.setFieldName(field)
	return t
}

// Rename sets the current Title value
func (t *Title) Rename(title string) error {
	memoized := t.value
	t.value = title
	if err := t.IsValid(); err != nil {
		t.value = memoized
		return err
	}

	return nil
}

// IsValid validates the current Title value(s)
func (t Title) IsValid() error {
	//	rules
	//	a.	max length 512 characters
	if len(t.value) > 512 {
		return exception.NewFieldRange(t.fieldName, "1", "512")
	}

	return nil
}

//	--	UTILS	--

// setFieldName sanitizes the given field name, if field empty then sets "title" by default
func (t *Title) setFieldName(f string) {
	if f == "" {
		t.fieldName = "title"
		return
	}

	t.fieldName = f
}

//	--	PRIMITIVES	--

// String returns the current Title value
func (t Title) String() string {
	return t.value
}

package value

import (
	"github.com/alexandria-oss/common-go/exception"
)

// Description text formatted using specific standards used to describe an aggregate/entity
type Description struct {
	value string

	fieldName string
}

// NewDescription creates a valid Description
func NewDescription(field, description string) (*Description, error) {
	d := new(Description)
	d.setFieldName(field)
	if err := d.Change(description); err != nil {
		return nil, err
	}

	return d, nil
}

// NewDescriptionFromPrimitive creates a Description without validating for marshaling purposes
func NewDescriptionFromPrimitive(field, description string) *Description {
	d := &Description{
		value: description,
	}
	d.setFieldName(field)
	return d
}

// Change sets the current Description value
func (d *Description) Change(description string) error {
	memoized := d.value
	d.value = description
	if err := d.IsValid(); err != nil {
		d.value = memoized
		return err
	}

	return nil
}

// IsValid validates the current Description value(s)
func (d Description) IsValid() error {
	//	rules
	//	a.	max length 1024 characters
	if len(d.value) > 1024 {
		return exception.NewFieldRange(d.fieldName, "1", "1024")
	}

	return nil
}

//	--	UTILS	--

// setFieldName sanitizes the given field name, if field empty then sets "description" by default
func (d *Description) setFieldName(f string) {
	if f == "" {
		d.fieldName = "description"
		return
	}

	d.fieldName = f
}

//	--	PRIMITIVES	--

// String returns the current Description value
func (d Description) String() string {
	return d.value
}

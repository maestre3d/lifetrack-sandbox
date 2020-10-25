package value

import (
	"strings"

	"github.com/alexandria-oss/common-go/exception"
)

// Image URL link to a remote bitmap file, likely to be stored in LifeTrack CDN bucket
type Image struct {
	value string

	fieldName string
}

// NewImage creates a valid Image
func NewImage(field, image string) (*Image, error) {
	i := &Image{
		value:     image,
		fieldName: "",
	}
	i.setFieldName(field)

	if err := i.IsValid(); err != nil {
		return nil, err
	}

	return i, nil
}

// NewImageFromPrimitive creates an Image without validating for marshaling purposes
func NewImageFromPrimitive(field, image string) *Image {
	i := &Image{value: image, fieldName: ""}
	i.setFieldName(field)
	return i
}

// Save stores the given image URL
func (i *Image) Save(image string) error {
	memoized := i.value
	i.value = image
	if err := i.IsValid(); err != nil {
		i.value = memoized
		return err
	}

	return nil
}

// IsValid validates the current Image value(s)
func (i Image) IsValid() error {
	//	rules
	//	a.	https-only
	//	b.	max range 2048 characters
	//	c.	.jpeg, .jpg, .png and .webp image formats only

	// business boolean expression(s)
	validImage := !strings.HasSuffix(i.value, ".jpeg") && !strings.HasSuffix(i.value, ".jpg") &&
		!strings.HasSuffix(i.value, ".png") && !strings.HasSuffix(i.value, ".webp")

	switch {
	case i.value != "" && !strings.HasPrefix(strings.ToLower(i.value), "https://"):
		return exception.NewFieldFormat(i.fieldName, "https link")
	case len(i.value) > 2048:
		return exception.NewFieldRange(i.fieldName, "1", "2048")
	case i.value != "" && validImage:
		return exception.NewFieldFormat(i.fieldName, "[jpeg, jpg, png, webp)")
	}

	return nil
}

//	--	UTILS	--

// setFieldName sanitizes the given field name, if field empty then sets "image" by default
func (i *Image) setFieldName(f string) {
	if f == "" {
		i.fieldName = "image"
		return
	}

	i.fieldName = strings.ToLower(f)
}

//	--	PRIMITIVES	--

// String returns the current Image value
func (i Image) String() string {
	return i.value
}

package value

import (
	gonanoid "github.com/matoous/go-nanoid"
)

// ID unique identifier
type ID struct {
	value string
}

// NewID creates an ID
func NewID() *ID {
	id := new(ID)
	id.Generate()
	return id
}

// NewIDFromPrimitive creates an ID from primitive without validation for marshaling purposes only
func NewIDFromPrimitive(id string) *ID {
	return &ID{value: id}
}

// Change updates the current ID value(s)
func (i *ID) Change(id string) error {
	memoized := i.value
	i.value = id
	if err := i.IsValid(); err != nil {
		i.value = memoized
		return err
	}

	return nil
}

// Generate generates and assigns a new unique identifier
func (i *ID) Generate() {
	id, _ := gonanoid.Nanoid()
	i.value = id
}

// IsValid validates the current ID value(s)
func (i ID) IsValid() error {
	/*
		if _, err := uuid.Parse(i.value); err != nil {
			return exception.NewFieldFormat("id", "uuid v4")
		}*/

	return nil
}

// FromPrimitive sets ID from given primitive for marshaling purposes only
func (i *ID) FromPrimitive(id string) {
	i.value = id
}

// String gets the current ID value
func (i ID) String() string {
	return i.value
}

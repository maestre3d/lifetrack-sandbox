package value

// Category is an aggregate.Category unique identifier value object
type Category struct {
	id string
}

// NewCategory creates a new aggregate.Category reference by ID
func NewCategory(id string) *Category {
	return &Category{id: id}
}

// IsEmpty checks if current Category value is empty
func (c Category) IsEmpty() bool {
	return c.id == ""
}

//	--	PRIMITIVES	--

// FromPrimitive sets the current Category value avoiding any domain validation for marshalling purposes
func (c *Category) FromPrimitive(id string) {
	c.id = id
}

// String returns current Category primitive value
func (c Category) String() string {
	return c.id
}

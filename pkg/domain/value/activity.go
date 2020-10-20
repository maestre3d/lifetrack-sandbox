package value

// Activity is an aggregate.Activity unique identifier value object
type Activity struct {
	id string
}

// NewActivity creates a new aggregate.Activity reference by ID
func NewActivity(id string) *Activity {
	return &Activity{id: id}
}

// IsEmpty checks if current Activity value is empty
func (a Activity) IsEmpty() bool {
	return a.id == ""
}

//	--	PRIMITIVES	--

// FromPrimitive sets the current Activity value avoiding any domain validation for marshalling purposes
func (a *Activity) FromPrimitive(id string) {
	a.id = id
}

// String returns current Activity primitive value
func (a Activity) String() string {
	return a.id
}

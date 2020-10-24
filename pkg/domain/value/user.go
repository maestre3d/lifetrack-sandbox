package value

// User is an aggregate.User unique identifier value object
type User struct {
	id string
}

// NewUser creates a new aggregate.User reference by ID
func NewUser(id string) *User {
	return &User{id: id}
}

// IsEmpty checks if current User value is empty
func (u User) IsEmpty() bool {
	return u.id == ""
}

//	--	PRIMITIVES	--

// FromPrimitive sets the current User value avoiding any domain validation for marshalling purposes
func (u *User) FromPrimitive(id string) {
	u.id = id
}

// String returns current User primitive value
func (u User) String() string {
	return u.id
}

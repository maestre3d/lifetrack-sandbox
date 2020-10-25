package model

// Category aggregate.Category primitive model
type Category struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	TargetTime  int64  `json:"target_time,omitempty"`
	Picture     string `json:"picture,omitempty"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	Active      bool   `json:"active"`
}

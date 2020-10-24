package model

// Activity aggregate.Activity primitive model
type Activity struct {
	ID         string `json:"id"`
	CategoryID string `json:"category_id"`
	Title      string `json:"title"`
	Picture    string `json:"picture"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
	Active     bool   `json:"active"`
}

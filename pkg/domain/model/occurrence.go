package model

// Occurrence is the aggregate.Occurrence primitive model
type Occurrence struct {
	ID            string `json:"id"`
	ActivityID    string `json:"activity_id"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	TotalDuration int64  `json:"total_duration"`
	CreateTime    int64  `json:"create_time"`
	UpdateTime    int64  `json:"update_time"`
}

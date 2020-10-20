package event

import "encoding/json"

// Domain is an action that has happened inside the ecosystem business contexts
type Domain struct {
	CorrelationID string      `json:"correlation_id"`
	Topic         string      `json:"topic"`
	Publisher     string      `json:"publisher"`
	Action        string      `json:"action"`
	PublishTime   int64       `json:"publish_time"`
	AggregateID   string      `json:"aggregate_id"`
	Body          interface{} `json:"body,omitempty"`
}

// MarshalBinary parses current Domain event into a JSON binary
func (d Domain) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

// UnmarshalBinary parses JSON binary into the current Domain event
func (d *Domain) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

package event

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/exceptions"

	"github.com/google/uuid"
)

// Domain is an action that has happened inside the ecosystem business contexts
type Domain struct {
	// CorrelationID distributed unique identifier for events/transactions
	CorrelationID string `json:"correlation_id"`
	// Publisher service or serverless function withing LifeTrack ecosystem that triggered the event
	Publisher string `json:"publisher"`
	// Action operation that has occurred
	Action string `json:"action"`
	// Topic event unique name, the event will be published and consumed using it
	Topic string `json:"topic"`
	// PublishTime UNIX timestamp the event was published
	PublishTime int64 `json:"publish_time"`
	// AggregateID unique identifier of the aggregate that throw the event
	AggregateID string `json:"aggregate_id"`
	// Body message of the event, if side-effect then it should contain a model of an aggregate
	Body interface{} `json:"body,omitempty"`
	// Version application/function version
	//	if empty, value will be set to "1.0.0"
	Version string `json:"version"`
	// Stage application/function development stage
	//	if empty, value will be set to "prod"
	Stage string `json:"stage"`
}

// DomainArgs arguments required to create a domain event
type DomainArgs struct {
	// Caller module or domain context who did the operation
	//	e.g. occurrence, activity, category
	Caller string
	// Action operation that has occurred
	Action string
	// AggregateName aggregate's name from domain
	AggregateName string
	// AggregateID unique identifier of the aggregate that throw the event
	AggregateID string
	// Body message of the event, if side-effect then it should contain a model of an aggregate
	//	value can be nil
	Body interface{}
}

// NewDomainEvent creates a new Domain event with formatted fields
func NewDomainEvent(args DomainArgs) *Domain {
	// Default topic format: "ecosystem_name.domain_context.aggregate_name.action"
	// (e.g. "lifetrack.tracker.activity.created")
	return &Domain{
		CorrelationID: uuid.New().String(),
		Publisher: fmt.Sprintf("lifetrack.%s.%s", strings.ToLower(args.Caller),
			strings.ToLower(args.AggregateName)),
		Action: strings.ToLower(args.Action),
		Topic: fmt.Sprintf("%s.%s.%s", strings.ToLower(args.Caller), strings.ToLower(args.AggregateName),
			strings.ToLower(args.Action)),
		PublishTime: time.Now().UTC().Unix(),
		AggregateID: args.AggregateID,
		Body:        args.Body,
		Version:     "1.0.0",
		Stage:       "prod",
	}
}

// MarshalBinary parses current Domain event into a JSON binary
func (d Domain) MarshalBinary() ([]byte, error) {
	j, err := json.Marshal(d)
	if err != nil {
		return nil, exceptions.ErrInvalidJSONEvent
	}

	return j, nil
}

// UnmarshalBinary parses JSON binary into the current Domain event
func (d *Domain) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, d); err != nil {
		return exceptions.ErrInvalidEvent
	}

	return nil
}

// TopicUnderscored retrieves a topic name with underscore format
func (d Domain) TopicUnderscored() string {
	return strings.Replace(strings.ToUpper(d.Topic), ".", "_", -1)
}

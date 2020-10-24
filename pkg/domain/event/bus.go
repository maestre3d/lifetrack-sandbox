package event

import "context"

// 	Bus connects a module to the LifeTrack ecosystem through domain and
//	integration events
type Bus interface {
	// Publish sends a group of Domain events to all subscribers
	Publish(ctx context.Context, e ...Domain) error
	// SubscribeTo adds a new subscription to an specific topic, returns a channel of Domain if exists
	SubscribeTo(ctx context.Context, topic string) (chan Domain, error)
}

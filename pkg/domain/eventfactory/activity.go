package eventfactory

import (
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/event"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"
)

type Activity struct{}

var activityName = "activity"

// Added triggered when a new aggregate.Activity has been added to a Category
func (a Activity) Added(activity model.Activity) event.Domain {
	return *event.NewDomainEvent(event.DomainArgs{
		Caller:        tracker,
		AggregateName: activityName,
		Action:        create,
		AggregateID:   activity.ID,
		Body:          activity,
	})
}

// Updated triggered when an existing aggregate.Activity has been updated
func (a Activity) Updated(activity model.Activity) event.Domain {
	return *event.NewDomainEvent(event.DomainArgs{
		Caller:        tracker,
		AggregateName: activityName,
		Action:        update,
		AggregateID:   activity.ID,
		Body:          activity,
	})
}

// Restored triggered when an existing aggregate.Activity has been restored
func (a Activity) Restored(activityID string) event.Domain {
	return *event.NewDomainEvent(event.DomainArgs{
		Caller:        tracker,
		AggregateName: activityName,
		Action:        restore,
		AggregateID:   activityID,
		Body:          nil,
	})
}

// Removed triggered when an existing aggregate.Activity has been removed softly
func (a Activity) Removed(activityID string) event.Domain {
	return *event.NewDomainEvent(event.DomainArgs{
		Caller:        tracker,
		AggregateName: activityName,
		Action:        remove,
		AggregateID:   activityID,
		Body:          nil,
	})
}

// HardRemoved triggered when an existing aggregate.Activity has been removed permanently
func (a Activity) HardRemoved(activityID string) event.Domain {
	return *event.NewDomainEvent(event.DomainArgs{
		Caller:        tracker,
		AggregateName: activityName,
		Action:        hardRemove,
		AggregateID:   activityID,
		Body:          nil,
	})
}

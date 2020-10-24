package eventfactory

import (
	"time"

	"github.com/google/uuid"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/event"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"
)

type Activity struct{}

// Added triggered when a new aggregate.Activity has been added to a Category
func (a Activity) Added(activity model.Activity) event.Domain {
	return event.Domain{
		CorrelationID: uuid.New().String(),
		Topic:         "activity_added",
		Publisher:     "lifetrack.activity",
		Action:        create,
		PublishTime:   time.Now().UTC().Unix(),
		AggregateID:   activity.ID,
		Body:          activity,
	}
}

// Updated triggered when an existing aggregate.Activity has been updated
func (a Activity) Updated(activity model.Activity) event.Domain {
	return event.Domain{
		CorrelationID: uuid.New().String(),
		Topic:         "activity_updated",
		Publisher:     "lifetrack.activity",
		Action:        update,
		PublishTime:   time.Now().UTC().Unix(),
		AggregateID:   activity.ID,
		Body:          activity,
	}
}

// Restored triggered when an existing aggregate.Activity has been restored
func (a Activity) Restored(activityID string) event.Domain {
	return event.Domain{
		CorrelationID: uuid.New().String(),
		Topic:         "activity_restored",
		Publisher:     "lifetrack.activity",
		Action:        restore,
		PublishTime:   time.Now().UTC().Unix(),
		AggregateID:   activityID,
		Body:          nil,
	}
}

// Removed triggered when an existing aggregate.Activity has been removed softly
func (a Activity) Removed(activityID string) event.Domain {
	return event.Domain{
		CorrelationID: uuid.New().String(),
		Topic:         "activity_removed",
		Publisher:     "lifetrack.activity",
		Action:        remove,
		PublishTime:   time.Now().UTC().Unix(),
		AggregateID:   activityID,
		Body:          nil,
	}
}

// HardRemoved triggered when an existing aggregate.Activity has been removed permanently
func (a Activity) HardRemoved(activityID string) event.Domain {
	return event.Domain{
		CorrelationID: uuid.New().String(),
		Topic:         "activity_hard_removed",
		Publisher:     "lifetrack.activity",
		Action:        hardRemove,
		PublishTime:   time.Now().UTC().Unix(),
		AggregateID:   activityID,
		Body:          nil,
	}
}

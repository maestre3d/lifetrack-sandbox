package eventfactory

import (
	"time"

	"github.com/google/uuid"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/event"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"
)

type Occurrence struct{}

// ActivityOccurred triggered when a new Occurrence has been added to an Activity
func (o Occurrence) ActivityOccurred(occurrence model.Occurrence) event.Domain {
	return event.Domain{
		CorrelationID: uuid.New().String(),
		Topic:         "activity_occurred",
		Publisher:     "lifetrack.occurrence",
		Action:        create,
		PublishTime:   time.Now().UTC().Unix(),
		AggregateID:   occurrence.ID,
		Body:          occurrence,
	}
}

// Updated triggered when an Occurrence has been updated
func (o Occurrence) Updated(occurrence model.Occurrence) event.Domain {
	return event.Domain{
		CorrelationID: uuid.New().String(),
		Topic:         "occurrence_updated",
		Publisher:     "lifetrack.occurrence",
		Action:        update,
		PublishTime:   time.Now().UTC().Unix(),
		AggregateID:   occurrence.ID,
		Body:          occurrence,
	}
}

// HardRemoved triggered when an Occurrence has been removed permanently
func (o Occurrence) HardRemoved(occurrenceID string) event.Domain {
	return event.Domain{
		CorrelationID: uuid.New().String(),
		Topic:         "occurrence_hard_removed",
		Publisher:     "lifetrack.occurrence",
		Action:        hardRemove,
		PublishTime:   time.Now().UTC().Unix(),
		AggregateID:   occurrenceID,
		Body:          nil,
	}
}

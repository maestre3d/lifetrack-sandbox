package eventfactory

import (
	"time"

	"github.com/google/uuid"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/event"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"
)

type Occurrence struct{}

// ActivityOccurred is triggered when a new Occurrence has been added to an activity
func (o Occurrence) ActivityOccurred(occurrence model.Occurrence) event.Domain {
	return event.Domain{
		CorrelationID: uuid.New().String(),
		Topic:         "activity_occurred",
		Publisher:     "lifetrack.occurrence",
		Action:        "create",
		PublishTime:   time.Now().UTC().Unix(),
		AggregateID:   occurrence.ID,
		Body:          occurrence,
	}
}

// Updated is triggered when an Occurrence has been updated
func (o Occurrence) Updated(occurrence model.Occurrence) event.Domain {
	return event.Domain{
		CorrelationID: uuid.New().String(),
		Topic:         "occurrence_updated",
		Publisher:     "lifetrack.occurrence",
		Action:        "update",
		PublishTime:   time.Now().UTC().Unix(),
		AggregateID:   occurrence.ID,
		Body:          occurrence,
	}
}

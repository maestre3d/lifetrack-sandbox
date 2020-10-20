package persistence

import (
	"context"
	"sync"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/exceptions"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
)

type OccurrenceInMemory struct {
	db map[string]*aggregate.Occurrence
	mu *sync.RWMutex
}

func NewOccurrenceInMemory() *OccurrenceInMemory {
	return &OccurrenceInMemory{db: map[string]*aggregate.Occurrence{}, mu: new(sync.RWMutex)}
}

func (o *OccurrenceInMemory) Save(_ context.Context, occurrence aggregate.Occurrence) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.db[occurrence.ID()] = &occurrence
	return nil
}

func (o *OccurrenceInMemory) Fetch(_ context.Context, id string) (*aggregate.Occurrence, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	oc := o.db[id]
	if oc == nil {
		return nil, exceptions.ErrOccurrenceNotFound
	}

	return oc, nil
}

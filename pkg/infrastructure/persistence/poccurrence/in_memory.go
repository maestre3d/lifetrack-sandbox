package poccurrence

import (
	"context"
	"sync"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/exceptions"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

var (
	occurrenceInMem     *OccurrenceInMemory
	occurrenceInMemOnce = new(sync.Once)
)

type OccurrenceInMemory struct {
	db map[string]*aggregate.Occurrence
	mu *sync.RWMutex
}

func NewOccurrenceInMemory() *OccurrenceInMemory {
	// Singleton
	occurrenceInMemOnce.Do(func() {
		if occurrenceInMem == nil {
			occurrenceInMem = &OccurrenceInMemory{db: map[string]*aggregate.Occurrence{}, mu: new(sync.RWMutex)}
		}
	})
	return occurrenceInMem
}

func (o *OccurrenceInMemory) Save(_ context.Context, occurrence aggregate.Occurrence) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.db[occurrence.ID()] = &occurrence
	return nil
}

func (o *OccurrenceInMemory) Fetch(ctx context.Context, criteria repository.OccurrenceCriteria) ([]*aggregate.Occurrence, string, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if criteria.Limit == 0 {
		criteria.Limit = 100
	}

	fetchStrategy := o.setFetchStrategy(criteria)
	if fetchStrategy == nil {
		return nil, "", exceptions.ErrInvalidOccurrenceFilter
	}

	ocs, nextPage, err := fetchStrategy.Do(ctx, criteria)
	if len(ocs) == 0 && err == nil {
		return nil, "", exceptions.ErrOccurrenceNotFound
	}
	return ocs, nextPage, err
}

func (o *OccurrenceInMemory) Remove(_ context.Context, id string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.db[id] == nil {
		return exceptions.ErrOccurrenceNotFound
	}

	delete(o.db, id)
	return nil
}

// setFetchStrategy chooses a fetching strategy depending on criteria values, if none returns nil
func (o OccurrenceInMemory) setFetchStrategy(criteria repository.OccurrenceCriteria) FetchStrategy {
	switch {
	case criteria.ID != "":
		return FetchIDInMemory{db: o.db}
	case criteria.Activity != "":
		return FetchActivityInMemory{db: o.db}
	case criteria.Limit > 0 || criteria.Token != "":
		return FetchAllInMemory{db: o.db}
	default:
		return nil
	}
}

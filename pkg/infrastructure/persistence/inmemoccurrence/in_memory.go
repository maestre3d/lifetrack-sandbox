package inmemoccurrence

import (
	"context"
	"sync"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/exceptions"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

var (
	inMemorySingleton *InMemory
	inMemoryLock      = new(sync.Once)
)

type inMemoryDatabase map[string]*aggregate.Occurrence

type InMemory struct {
	db inMemoryDatabase
	mu *sync.RWMutex
}

func NewInMemory() *InMemory {
	inMemoryLock.Do(func() {
		inMemorySingleton = &InMemory{db: map[string]*aggregate.Occurrence{}, mu: new(sync.RWMutex)}
	})
	return inMemorySingleton
}

func (r *InMemory) Save(_ context.Context, occurrence aggregate.Occurrence) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.db[occurrence.ID()] = &occurrence
	return nil
}

func (r *InMemory) Fetch(ctx context.Context, criteria repository.OccurrenceCriteria) ([]*aggregate.Occurrence, string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ocs, nextPage, err := r.setFetchStrategy(criteria).Do(ctx, criteria)
	if len(ocs) == 0 && err == nil {
		return nil, "", exceptions.ErrOccurrenceNotFound
	}
	return ocs, nextPage, err
}

func (r *InMemory) Remove(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.db[id] == nil {
		return exceptions.ErrOccurrenceNotFound
	}

	delete(r.db, id)
	return nil
}

// setFetchStrategy chooses a fetching strategy depending on criteria values, if none returns nil
func (r InMemory) setFetchStrategy(criteria repository.OccurrenceCriteria) fetchStrategy {
	switch {
	case criteria.ID != "":
		return fetchIDInMemory{db: r.db}
	case criteria.Activity != "":
		return fetchActivityInMemory{db: r.db}
	default:
		return fetchAllInMemory{db: r.db}
	}
}

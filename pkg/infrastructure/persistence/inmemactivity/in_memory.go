package inmemactivity

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

type inMemoryDatabase map[string]*aggregate.Activity

type InMemory struct {
	db inMemoryDatabase
	mu *sync.RWMutex
}

func NewInMemory() *InMemory {
	inMemoryLock.Do(func() {
		inMemorySingleton = &InMemory{db: map[string]*aggregate.Activity{}, mu: new(sync.RWMutex)}
	})
	return inMemorySingleton
}

func (r *InMemory) Save(_ context.Context, activity aggregate.Activity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.db[activity.ID()] = &activity
	return nil
}

func (r *InMemory) Fetch(ctx context.Context, criteria repository.ActivityCriteria) ([]*aggregate.Activity, string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	acts, nextPage, err := r.setFetchStrategy(criteria).Do(ctx, criteria)
	if len(acts) == 0 && err == nil {
		return nil, "", exceptions.ErrActivityNotFound
	}
	return acts, nextPage, err
}

func (r *InMemory) Remove(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.db[id] == nil {
		return exceptions.ErrActivityNotFound
	}

	delete(r.db, id)
	return nil
}

// setFetchStrategy chooses a fetching strategy depending on criteria values, if none returns nil
func (r InMemory) setFetchStrategy(criteria repository.ActivityCriteria) fetchStrategy {
	switch {
	case criteria.ID != "":
		return fetchIDInMemory{db: r.db}
	case criteria.Title != "":
		return fetchTitleInMemory{db: r.db}
	case criteria.Category != "":
		return fetchCategoryInMemory{db: r.db}
	default:
		return fetchAllInMemory{db: r.db}
	}
}

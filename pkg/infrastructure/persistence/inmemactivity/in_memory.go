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

func (o *InMemory) Save(_ context.Context, activity aggregate.Activity) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.db[activity.ID()] = &activity
	return nil
}

func (o *InMemory) Fetch(ctx context.Context, criteria repository.ActivityCriteria) ([]*aggregate.Activity, string, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if criteria.Limit == 0 {
		criteria.Limit = 100
	}

	fetchStrategy := o.setFetchStrategy(criteria)
	if fetchStrategy == nil {
		return nil, "", exceptions.ErrInvalidActivityFilter
	}

	ocs, nextPage, err := fetchStrategy.Do(ctx, criteria)
	if len(ocs) == 0 && err == nil {
		return nil, "", exceptions.ErrActivityNotFound
	}
	return ocs, nextPage, err
}

func (o *InMemory) Remove(_ context.Context, id string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.db[id] == nil {
		return exceptions.ErrActivityNotFound
	}

	delete(o.db, id)
	return nil
}

// setFetchStrategy chooses a fetching strategy depending on criteria values, if none returns nil
func (o InMemory) setFetchStrategy(criteria repository.ActivityCriteria) fetchStrategy {
	switch {
	case criteria.ID != "":
		return fetchIDInMemory{db: o.db}
	case criteria.Title != "":
		return fetchTitleInMemory{db: o.db}
	case criteria.Category != "":
		return fetchCategoryInMemory{db: o.db}
	case criteria.Limit > 0 || criteria.Token != "":
		return fetchAllInMemory{db: o.db}
	default:
		return nil
	}
}

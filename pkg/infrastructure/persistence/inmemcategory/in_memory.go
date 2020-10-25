package inmemcategory

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

type inMemoryDatabase map[string]*aggregate.Category

type InMemory struct {
	db inMemoryDatabase
	mu *sync.RWMutex
}

func NewInMemory() *InMemory {
	inMemoryLock.Do(func() {
		inMemorySingleton = &InMemory{db: map[string]*aggregate.Category{}, mu: new(sync.RWMutex)}
	})
	return inMemorySingleton
}

func (r *InMemory) Save(_ context.Context, Category aggregate.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.db[Category.ID()] = &Category
	return nil
}

func (r *InMemory) Fetch(ctx context.Context, criteria repository.CategoryCriteria) ([]*aggregate.Category, string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cats, nextPage, err := r.setFetchStrategy(criteria).Do(ctx, criteria)
	if len(cats) == 0 && err == nil {
		return nil, "", exceptions.ErrCategoryNotFound
	}
	return cats, nextPage, err
}

func (r *InMemory) Remove(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.db[id] == nil {
		return exceptions.ErrCategoryNotFound
	}

	delete(r.db, id)
	return nil
}

// setFetchStrategy chooses a fetching strategy depending on criteria values, if none returns nil
func (r InMemory) setFetchStrategy(criteria repository.CategoryCriteria) fetchStrategy {
	switch {
	case criteria.ID != "":
		return fetchIDInMemory{db: r.db}
	case criteria.Name != "":
		return fetchNameInMemory{db: r.db}
	case criteria.Keyword != "":
		return fetchKeywordInMemory{db: r.db}
	case criteria.User != "":
		return fetchUserInMemory{db: r.db}
	default:
		return fetchAllInMemory{db: r.db}
	}
}

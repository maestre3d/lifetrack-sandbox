package inmemcategory

import (
	"context"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// fetchStrategy fetching strategy for category repositories
type fetchStrategy interface {
	Do(ctx context.Context, criteria repository.CategoryCriteria) ([]*aggregate.Category, string, error)
}

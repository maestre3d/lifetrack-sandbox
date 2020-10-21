package poccurrence

import (
	"context"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// FetchStrategy fetching strategy for occurrence repositories
type FetchStrategy interface {
	Do(ctx context.Context, criteria repository.OccurrenceCriteria) ([]*aggregate.Occurrence, string, error)
}

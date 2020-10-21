package repository

import (
	"context"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
)

// OccurrenceCriteria sets the Occurrence fetching strategy
type OccurrenceCriteria struct {
	ID       string
	Activity string
	Limit    int64
	Token    string
}

// Occurrence handles aggregate.Occurrence persistence
type Occurrence interface {
	Save(ctx context.Context, occurrence aggregate.Occurrence) error
	Fetch(ctx context.Context, criteria OccurrenceCriteria) ([]*aggregate.Occurrence, string, error)
	Remove(ctx context.Context, id string) error
}

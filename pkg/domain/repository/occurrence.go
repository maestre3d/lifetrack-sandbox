package repository

import (
	"context"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
)

// Occurrence handles aggregate.Occurrence persistence
type Occurrence interface {
	Save(ctx context.Context, occurrence aggregate.Occurrence) error
	Fetch(ctx context.Context, id string) (*aggregate.Occurrence, error)
}

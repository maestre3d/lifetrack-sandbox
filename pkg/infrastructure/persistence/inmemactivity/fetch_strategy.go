package inmemactivity

import (
	"context"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// fetchStrategy fetching strategy for activity repositories
type fetchStrategy interface {
	Do(ctx context.Context, criteria repository.ActivityCriteria) ([]*aggregate.Activity, string, error)
}

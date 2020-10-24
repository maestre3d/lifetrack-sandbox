package repository

import (
	"context"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
)

// ActivityCriteria sets the Activity fetching strategy
type ActivityCriteria struct {
	ID       string
	Activity string
	Limit    int64
	Token    string
}

// Activity handles aggregate.Activity persistence
type Activity interface {
	Save(ctx context.Context, activity aggregate.Activity) error
	Fetch(ctx context.Context, criteria ActivityCriteria) ([]*aggregate.Activity, string, error)
	Remove(ctx context.Context, id string) error
}

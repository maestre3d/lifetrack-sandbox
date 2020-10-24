package repository

import (
	"context"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
)

// CategoryCriteria sets the Category fetching strategy
type CategoryCriteria struct {
	ID   string
	User string
	Name string
	// Keyword fetch Category containing the given keyword in name and description fields
	Keyword string
	Limit   int64
	Token   string
}

// Category handles aggregate.Category persistence
type Category interface {
	Save(ctx context.Context, activity aggregate.Category) error
	Fetch(ctx context.Context, criteria CategoryCriteria) ([]*aggregate.Category, string, error)
	Remove(ctx context.Context, id string) error
}

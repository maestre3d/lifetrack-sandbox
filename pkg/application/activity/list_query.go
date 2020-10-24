package activity

import (
	"context"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/adapter"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// ListQuery requests a set of model.Activity
type ListQuery struct {
	repo repository.Activity
}

// Filter sets Activity fetching strategy
//	anti-corruption struct
type Filter struct {
	CategoryID string `json:"category_id"`
	Title      string `json:"title"`
	Limit      int64  `json:"limit"`
	Token      string `json:"token"`
}

// NewListQuery creates a ListQuery
func NewListQuery(r repository.Activity) *ListQuery {
	return &ListQuery{repo: r}
}

func (q ListQuery) Query(ctx context.Context, filter Filter) ([]*model.Activity, string, error) {
	if filter.Limit == 0 {
		filter.Limit = 100
	}

	activities, nextToken, err := q.repo.Fetch(ctx, repository.ActivityCriteria{
		Category: filter.CategoryID,
		Title:    filter.Title,
		Limit:    filter.Limit,
		Token:    filter.Token,
	})
	return adapter.BulkUnmarshalPrimitiveActivity(activities), nextToken, err
}

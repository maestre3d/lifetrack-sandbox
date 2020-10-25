package category

import (
	"context"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/adapter"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// ListQuery requests a set of model.Category
type ListQuery struct {
	repo repository.Category
}

// Filter sets Category fetching strategy
//	anti-corruption struct
type Filter struct {
	UserID  string `json:"user_id"`
	Name    string `json:"name"`
	Keyword string `json:"keyword"`
	Limit   int64  `json:"limit"`
	Token   string `json:"token"`
}

// NewListQuery creates a ListQuery
func NewListQuery(r repository.Category) *ListQuery {
	return &ListQuery{repo: r}
}

func (q ListQuery) Query(ctx context.Context, filter Filter) ([]*model.Category, string, error) {
	if filter.Limit == 0 {
		filter.Limit = 100
	}

	categories, nextToken, err := q.repo.Fetch(ctx, repository.CategoryCriteria{
		User:    filter.UserID,
		Name:    filter.Name,
		Keyword: filter.Keyword,
		Limit:   filter.Limit,
		Token:   filter.Token,
	})
	return adapter.BulkUnmarshalPrimitiveCategory(categories), nextToken, err
}

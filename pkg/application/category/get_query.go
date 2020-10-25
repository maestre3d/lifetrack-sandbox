package category

import (
	"context"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// GetQuery requests a single model.Category
type GetQuery struct {
	repo repository.Category
}

// NewGetQuery creates a GetQuery
func NewGetQuery(r repository.Category) *GetQuery {
	return &GetQuery{repo: r}
}

func (q GetQuery) Query(ctx context.Context, id string) (*model.Category, error) {
	cats, _, err := q.repo.Fetch(ctx, repository.CategoryCriteria{ID: id})
	if err != nil {
		return nil, err
	}

	return cats[0].MarshalPrimitive(), nil
}

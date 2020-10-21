package occurrence

import (
	"context"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// GetQuery requests a get Occurrence query
type GetQuery struct {
	repo repository.Occurrence
}

// NewGetQuery creates a new Get query
func NewGetQuery(r repository.Occurrence) *GetQuery {
	return &GetQuery{repo: r}
}

// Query handles Get requests
func (g GetQuery) Query(ctx context.Context, id string) (*model.Occurrence, error) {
	oc, _, err := g.repo.Fetch(ctx, repository.OccurrenceCriteria{ID: id})
	if err != nil {
		return nil, err
	}

	return oc[0].MarshalPrimitive(), nil
}

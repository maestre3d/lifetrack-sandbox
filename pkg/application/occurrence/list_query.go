package occurrence

import (
	"context"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/adapter"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// ListQuery requests a list of model.Occurrence
type ListQuery struct {
	repo repository.Occurrence
}

// NewListQuery creates a ListQuery
func NewListQuery(r repository.Occurrence) *ListQuery {
	return &ListQuery{repo: r}
}

// Filter sets the Occurrence fetching strategy
type Filter struct {
	Activity string `json:"activity"`
	Limit    int64  `json:"limit"`
	Token    string `json:"token"`
}

func (q ListQuery) Query(ctx context.Context, filter Filter) ([]*model.Occurrence, string, error) {
	ocs, nextPage, err := q.repo.Fetch(ctx, repository.OccurrenceCriteria{
		Activity: filter.Activity,
		Limit:    filter.Limit,
		Token:    filter.Token,
	})

	return adapter.BulkUnmarshalPrimitiveOccurrence(ocs), nextPage, err
}

package poccurrence

import (
	"context"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/exceptions"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// FetchIDInMemory strategy when criteria contains an ID
type FetchIDInMemory struct {
	db map[string]*aggregate.Occurrence
}

func (m FetchIDInMemory) Do(_ context.Context, criteria repository.OccurrenceCriteria) ([]*aggregate.Occurrence, string, error) {
	oc := m.db[criteria.ID]
	if oc == nil {
		return nil, "", exceptions.ErrOccurrenceNotFound
	}

	return []*aggregate.Occurrence{oc}, "", nil
}

// FetchAllInMemory strategy when criteria contains a limit above 0 or a page token
type FetchAllInMemory struct {
	db map[string]*aggregate.Occurrence
}

func (m FetchAllInMemory) Do(_ context.Context, criteria repository.OccurrenceCriteria) ([]*aggregate.Occurrence, string, error) {
	totalRows := criteria.Limit
	rows := make([]*aggregate.Occurrence, 0)
	for _, oc := range m.db {
		if totalRows == 0 {
			break
		}
		rows = append(rows, oc)
		totalRows--
	}

	return rows, "", nil
}

// FetchActivityInMemory strategy when criteria contains an Activity ID
type FetchActivityInMemory struct {
	db map[string]*aggregate.Occurrence
}

func (m FetchActivityInMemory) Do(_ context.Context, criteria repository.OccurrenceCriteria) ([]*aggregate.Occurrence, string, error) {
	totalRows := criteria.Limit
	rows := make([]*aggregate.Occurrence, 0)
	for _, oc := range m.db {
		if totalRows == 0 {
			break
		} else if oc.Activity().String() == criteria.Activity {
			rows = append(rows, oc)
			totalRows--
		}
	}

	return rows, "", nil
}

package poccurrence

import (
	"context"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/exceptions"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// fetchIDInMemory strategy when criteria contains an ID
type fetchIDInMemory struct {
	db map[string]*aggregate.Occurrence
}

func (m fetchIDInMemory) Do(_ context.Context, criteria repository.OccurrenceCriteria) ([]*aggregate.Occurrence, string, error) {
	oc := m.db[criteria.ID]
	if oc == nil {
		return nil, "", exceptions.ErrOccurrenceNotFound
	}

	return []*aggregate.Occurrence{oc}, "", nil
}

// fetchAllInMemory strategy when criteria contains a limit above 0 or a page token
type fetchAllInMemory struct {
	db map[string]*aggregate.Occurrence
}

func (m fetchAllInMemory) Do(_ context.Context, criteria repository.OccurrenceCriteria) ([]*aggregate.Occurrence, string, error) {
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

// fetchActivityInMemory strategy when criteria contains an Activity ID
type fetchActivityInMemory struct {
	db map[string]*aggregate.Occurrence
}

func (m fetchActivityInMemory) Do(_ context.Context, criteria repository.OccurrenceCriteria) ([]*aggregate.Occurrence, string, error) {
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

package inmemoccurrence

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
	rows := make([]*aggregate.Occurrence, 0)
	nextToken := ""
	for _, oc := range m.db {
		if criteria.Limit == 0 {
			nextToken = oc.ID()
			break
		}
		rows = append(rows, oc)
		criteria.Limit--
	}

	return rows, nextToken, nil
}

// fetchActivityInMemory strategy when criteria contains an Activity ID
type fetchActivityInMemory struct {
	db map[string]*aggregate.Occurrence
}

func (m fetchActivityInMemory) Do(_ context.Context, criteria repository.OccurrenceCriteria) ([]*aggregate.Occurrence, string, error) {
	rows := make([]*aggregate.Occurrence, 0)
	nextToken := ""
	for _, oc := range m.db {
		if criteria.Limit == 0 {
			nextToken = oc.ID()
			break
		} else if oc.Activity() == criteria.Activity {
			rows = append(rows, oc)
			criteria.Limit--
		}
	}

	return rows, nextToken, nil
}

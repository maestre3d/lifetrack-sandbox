package inmemactivity

import (
	"context"
	"strings"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/exceptions"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// fetchIDInMemory strategy when criteria contains an ID
type fetchIDInMemory struct {
	db map[string]*aggregate.Activity
}

func (m fetchIDInMemory) Do(_ context.Context, criteria repository.ActivityCriteria) ([]*aggregate.Activity, string, error) {
	act := m.db[criteria.ID]
	if act == nil {
		return nil, "", exceptions.ErrActivityNotFound
	}

	return []*aggregate.Activity{act}, "", nil
}

// fetchAllInMemory strategy when criteria contains a limit above 0 or a page token
type fetchAllInMemory struct {
	db map[string]*aggregate.Activity
}

func (m fetchAllInMemory) Do(_ context.Context, criteria repository.ActivityCriteria) ([]*aggregate.Activity, string, error) {
	totalRows := criteria.Limit
	rows := make([]*aggregate.Activity, 0)
	for _, act := range m.db {
		if totalRows == 0 {
			break
		}
		rows = append(rows, act)
		totalRows--
	}

	return rows, "", nil
}

// fetchTitleInMemory strategy when criteria contains a Category ID
type fetchTitleInMemory struct {
	db map[string]*aggregate.Activity
}

func (m fetchTitleInMemory) Do(_ context.Context, criteria repository.ActivityCriteria) ([]*aggregate.Activity, string, error) {
	totalRows := criteria.Limit
	rows := make([]*aggregate.Activity, 0)
	for _, act := range m.db {
		if totalRows == 0 {
			break
		} else if strings.Contains(act.Title(), criteria.Title) {
			rows = append(rows, act)
			totalRows--
		}
	}

	return rows, "", nil
}

// fetchCategoryInMemory strategy when criteria contains a Category ID
type fetchCategoryInMemory struct {
	db map[string]*aggregate.Activity
}

func (m fetchCategoryInMemory) Do(_ context.Context, criteria repository.ActivityCriteria) ([]*aggregate.Activity, string, error) {
	totalRows := criteria.Limit
	rows := make([]*aggregate.Activity, 0)
	for _, act := range m.db {
		if totalRows == 0 {
			break
		} else if act.Category() == criteria.Category {
			rows = append(rows, act)
			totalRows--
		}
	}

	return rows, "", nil
}

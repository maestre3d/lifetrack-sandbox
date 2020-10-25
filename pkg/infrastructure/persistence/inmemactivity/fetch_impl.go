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
	rows := make([]*aggregate.Activity, 0)
	nextToken := ""
	for _, act := range m.db {
		if criteria.Limit == 0 {
			nextToken = act.ID()
			break
		} else if criteria.Token != "" && criteria.Token == act.ID() {
			criteria.Token = ""
		}

		if criteria.Token == "" && act.State() {
			rows = append(rows, act)
			criteria.Limit--
		}
	}

	return rows, nextToken, nil
}

// fetchTitleInMemory strategy when criteria contains a title
type fetchTitleInMemory struct {
	db map[string]*aggregate.Activity
}

func (m fetchTitleInMemory) Do(_ context.Context, criteria repository.ActivityCriteria) ([]*aggregate.Activity, string, error) {
	rows := make([]*aggregate.Activity, 0)
	nextToken := ""
	for _, act := range m.db {
		if criteria.Limit == 0 {
			nextToken = act.ID()
			break
		} else if strings.Contains(strings.ToLower(act.Title()), strings.ToLower(criteria.Title)) &&
			act.State() {
			rows = append(rows, act)
			criteria.Limit--
		}
	}

	return rows, nextToken, nil
}

// fetchCategoryInMemory strategy when criteria contains a Category ID
type fetchCategoryInMemory struct {
	db map[string]*aggregate.Activity
}

func (m fetchCategoryInMemory) Do(_ context.Context, criteria repository.ActivityCriteria) ([]*aggregate.Activity, string, error) {
	rows := make([]*aggregate.Activity, 0)
	nextToken := ""
	for _, act := range m.db {
		if criteria.Limit == 0 {
			nextToken = act.ID()
			break
		} else if act.Category() == criteria.Category && act.State() {
			rows = append(rows, act)
			criteria.Limit--
		}
	}

	return rows, nextToken, nil
}

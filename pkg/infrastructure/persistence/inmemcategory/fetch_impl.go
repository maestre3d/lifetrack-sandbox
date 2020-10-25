package inmemcategory

import (
	"context"
	"strings"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/exceptions"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// fetchIDInMemory strategy when criteria contains an ID
type fetchIDInMemory struct {
	db map[string]*aggregate.Category
}

func (m fetchIDInMemory) Do(_ context.Context, criteria repository.CategoryCriteria) ([]*aggregate.Category, string, error) {
	act := m.db[criteria.ID]
	if act == nil {
		return nil, "", exceptions.ErrCategoryNotFound
	}

	return []*aggregate.Category{act}, "", nil
}

// fetchAllInMemory strategy when criteria contains a limit above 0 or a page token
type fetchAllInMemory struct {
	db map[string]*aggregate.Category
}

func (m fetchAllInMemory) Do(_ context.Context, criteria repository.CategoryCriteria) ([]*aggregate.Category, string, error) {
	rows := make([]*aggregate.Category, 0)
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

// fetchNameInMemory strategy when criteria contains a name
type fetchNameInMemory struct {
	db map[string]*aggregate.Category
}

func (m fetchNameInMemory) Do(_ context.Context, criteria repository.CategoryCriteria) ([]*aggregate.Category, string, error) {
	rows := make([]*aggregate.Category, 0)
	nextToken := ""
	for _, act := range m.db {
		if criteria.Limit == 0 {
			nextToken = act.ID()
			break
		} else if strings.Contains(strings.ToLower(act.Name()), strings.ToLower(criteria.Name)) &&
			act.State() {
			rows = append(rows, act)
			criteria.Limit--
		}
	}

	return rows, nextToken, nil
}

// fetchKeywordInMemory strategy when criteria contains a keyword
type fetchKeywordInMemory struct {
	db map[string]*aggregate.Category
}

func (m fetchKeywordInMemory) Do(_ context.Context, criteria repository.CategoryCriteria) ([]*aggregate.Category, string, error) {
	rows := make([]*aggregate.Category, 0)
	nextToken := ""
	for _, act := range m.db {
		if criteria.Limit == 0 {
			nextToken = act.ID()
			break
		}
		containsKeyword := strings.Contains(strings.ToLower(act.Name()), strings.ToLower(criteria.Keyword)) ||
			strings.Contains(strings.ToLower(act.Description()), strings.ToLower(criteria.Keyword))

		if (containsKeyword) && act.State() {
			rows = append(rows, act)
			criteria.Limit--
		}
	}

	return rows, nextToken, nil
}

// fetchUserInMemory strategy when criteria contains a Category ID
type fetchUserInMemory struct {
	db map[string]*aggregate.Category
}

func (m fetchUserInMemory) Do(_ context.Context, criteria repository.CategoryCriteria) ([]*aggregate.Category, string, error) {
	rows := make([]*aggregate.Category, 0)
	nextToken := ""
	for _, act := range m.db {
		if criteria.Limit == 0 {
			nextToken = act.ID()
			break
		} else if act.User() == criteria.User && act.State() {
			rows = append(rows, act)
			criteria.Limit--
		}
	}

	return rows, nextToken, nil
}

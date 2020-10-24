package adapter

import (
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"
)

// BulkUnmarshalPrimitiveCategory parses given aggregate.Category slice into a read model slice
func BulkUnmarshalPrimitiveCategory(categories []*aggregate.Category) []*model.Category {
	cats := make([]*model.Category, 0)
	for _, cat := range categories {
		cats = append(cats, cat.MarshalPrimitive())
	}

	return cats
}

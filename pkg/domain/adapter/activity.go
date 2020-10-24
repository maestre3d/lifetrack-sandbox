package adapter

import (
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"
)

// BulkUnmarshalPrimitiveActivity parses given aggregate.Activity slice into a read model slice
func BulkUnmarshalPrimitiveActivity(activities []*aggregate.Activity) []*model.Activity {
	acts := make([]*model.Activity, 0)
	for _, ac := range activities {
		acts = append(acts, ac.MarshalPrimitive())
	}

	return acts
}

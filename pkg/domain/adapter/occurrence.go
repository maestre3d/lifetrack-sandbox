package adapter

import (
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"
)

// BulkUnmarshalPrimitiveOccurrence parses given aggregate.Occurrence slice into a read model slice
func BulkUnmarshalPrimitiveOccurrence(occurrences []*aggregate.Occurrence) []*model.Occurrence {
	ocs := make([]*model.Occurrence, 0)
	for _, oc := range occurrences {
		ocs = append(ocs, oc.MarshalPrimitive())
	}

	return ocs
}

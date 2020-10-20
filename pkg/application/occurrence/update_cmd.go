package occurrence

import (
	"context"
	"time"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// UpdateCommand requests an update of an specific Occurrence
type UpdateCommand struct {
	Ctx context.Context
	ID  string

	StartTime int64
	EndTime   int64

	ActivityID string
}

// UpdateCommandHandler handles UpdateCommand requests
type UpdateCommandHandler struct {
	repo repository.Occurrence
}

// NewUpdateCommandHandler creates an UpdateCommandHandler
func NewUpdateCommandHandler(r repository.Occurrence) *UpdateCommandHandler {
	return &UpdateCommandHandler{repo: r}
}

// Invoke handle an UpdateCommand request
func (h UpdateCommandHandler) Invoke(cmd UpdateCommand) error {
	if cmd.ID == "" {
		return exception.NewRequiredField("occurrence_id")
	}

	oc, err := h.repo.Fetch(cmd.Ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err := h.updater(cmd, oc); err != nil {
		return err
	}

	// TODO: Add event bus, publish events to AWS SNS/SQS or EventBridge
	return h.repo.Save(cmd.Ctx, *oc)
}

// updater updates aggregate.Occurrence using strategies
func (h UpdateCommandHandler) updater(cmd UpdateCommand, oc *aggregate.Occurrence) error {
	if cmd.ActivityID != "" {
		if err := oc.ChangeActivity(cmd.ActivityID); err != nil {
			return err
		}
	}

	return oc.EditTimes(time.Unix(cmd.StartTime, 0), time.Unix(cmd.EndTime, 0))
}

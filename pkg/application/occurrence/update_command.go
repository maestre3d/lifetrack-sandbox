package occurrence

import (
	"context"
	"time"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/event"
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
	bus  event.Bus
}

// NewUpdateCommandHandler creates an UpdateCommandHandler
func NewUpdateCommandHandler(r repository.Occurrence, b event.Bus) *UpdateCommandHandler {
	return &UpdateCommandHandler{repo: r, bus: b}
}

// Invoke handle an UpdateCommand request
func (h UpdateCommandHandler) Invoke(cmd UpdateCommand) error {
	if cmd.ID == "" {
		return exception.NewRequiredField("occurrencecurrence_id")
	}

	occurrence, _, err := h.repo.Fetch(cmd.Ctx, repository.OccurrenceCriteria{ID: cmd.ID})
	if err != nil {
		return err
	}
	snapshot := *occurrence[0]

	if err = h.updater(cmd, occurrence[0]); err != nil {
		return err
	}

	return h.persist(cmd.Ctx, occurrence[0], snapshot)
}

// updater updates aggregate.Occurrence using strategies
func (h UpdateCommandHandler) updater(cmd UpdateCommand, occurrence *aggregate.Occurrence) error {
	isCmdEmpty := cmd.ActivityID == "" && cmd.StartTime == 0 && cmd.EndTime == 0
	if cmd.ActivityID != "" {
		if err := occurrence.ChangeActivity(cmd.ActivityID); err != nil {
			return err
		}
	} else if isCmdEmpty {
		return exception.NewRequiredField("activity_id or times (start_time, end_time)")
	}

	return occurrence.EditTimes(time.Unix(cmd.StartTime, 0), time.Unix(cmd.EndTime, 0))
}

func (h UpdateCommandHandler) persist(ctx context.Context, occurrence *aggregate.Occurrence,
	snapshot aggregate.Occurrence) error {
	if err := h.repo.Save(ctx, *occurrence); err != nil {
		return err
	}

	return h.pushEvents(ctx, occurrence, snapshot)
}

func (h UpdateCommandHandler) pushEvents(ctx context.Context, occurrence *aggregate.Occurrence,
	snapshot aggregate.Occurrence) error {
	if err := h.bus.Publish(ctx, occurrence.PullEvents()...); err != nil {
		// rollback
		if errR := h.repo.Save(ctx, snapshot); errR != nil {
			return errR
		}
		return err
	}

	return nil
}

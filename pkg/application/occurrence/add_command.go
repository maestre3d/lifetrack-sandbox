package occurrence

import (
	"context"
	"time"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/event"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// AddCommand requests an Occurrence creation
type AddCommand struct {
	Ctx        context.Context
	ActivityID string
	StartTime  int64
	EndTime    int64
}

// AddCommandHandler handles AddCommand(s)
type AddCommandHandler struct {
	repo repository.Occurrence
	bus  event.Bus
}

// NewAddCommandHandler Adds a AddCommandHandler
func NewAddCommandHandler(r repository.Occurrence, b event.Bus) *AddCommandHandler {
	return &AddCommandHandler{repo: r, bus: b}
}

// Invoke handle a AddCommand request, returns Occurrence ID if succeed
func (h AddCommandHandler) Invoke(cmd AddCommand) (string, error) {
	occurrence, err := h.creator(cmd.ActivityID, cmd.StartTime, cmd.EndTime)
	if err != nil {
		return "", err
	}

	if err = h.persist(cmd.Ctx, occurrence); err != nil {
		return "", err
	}
	return occurrence.ID(), nil
}

// creator Adds an aggregate.Occurrence parsing required data
func (h AddCommandHandler) creator(activity string, startTime, endTime int64) (*aggregate.Occurrence, error) {
	return aggregate.NewOccurrence(activity, time.Unix(startTime, 0), time.Unix(endTime, 0))
}

func (h AddCommandHandler) persist(ctx context.Context, oc *aggregate.Occurrence) error {
	if err := h.repo.Save(ctx, *oc); err != nil {
		return err
	}
	return h.pushEvents(ctx, oc)
}

func (h AddCommandHandler) pushEvents(ctx context.Context, oc *aggregate.Occurrence) error {
	if err := h.bus.Publish(ctx, oc.PullEvents()...); err != nil {
		// rollback
		if errR := h.repo.Remove(ctx, oc.ID()); errR != nil {
			return errR
		}
		return err
	}

	return nil
}

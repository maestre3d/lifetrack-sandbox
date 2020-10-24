package occurrence

import (
	"context"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/event"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/eventfactory"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// RemoveCommand requests a Occurrence removal
type RemoveCommand struct {
	Ctx context.Context
	ID  string
}

// RemoveCommandHandler handles RemoveCommand(s)
type RemoveCommandHandler struct {
	repo repository.Occurrence
	bus  event.Bus
}

// NewRemoveCommandHandler creates a RemoveCommandHandler
func NewRemoveCommandHandler(r repository.Occurrence, b event.Bus) *RemoveCommandHandler {
	return &RemoveCommandHandler{repo: r, bus: b}
}

func (h RemoveCommandHandler) Invoke(cmd RemoveCommand) error {
	if cmd.ID == "" {
		return exception.NewRequiredField("occurrence_id")
	}
	return h.persist(cmd.Ctx, cmd.ID)
}

func (h RemoveCommandHandler) persist(ctx context.Context, id string) error {
	// required to verify entity exists and for rollback ops
	ocs, _, err := h.repo.Fetch(ctx, repository.OccurrenceCriteria{ID: id})
	if err != nil {
		return err
	}

	if err = h.repo.Remove(ctx, ocs[0].ID()); err != nil {
		return err
	}

	return h.pushEvents(ctx, ocs[0])
}

func (h RemoveCommandHandler) pushEvents(ctx context.Context, oc *aggregate.Occurrence) error {
	if err := h.bus.Publish(ctx, eventfactory.Occurrence{}.HardRemoved(oc.ID())); err != nil {
		//  rollback
		if errR := h.repo.Save(ctx, *oc); errR != nil {
			return errR
		}
		return err
	}

	return nil
}

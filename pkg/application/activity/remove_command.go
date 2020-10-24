package activity

import (
	"context"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/event"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/eventfactory"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// RemoveCommand requests a permanent removal of an aggregate.Activity
type RemoveCommand struct {
	Ctx context.Context
	ID  string
}

// RemoveCommandHandler handles RemoveCommand request(s)
type RemoveCommandHandler struct {
	repo repository.Activity
	bus  event.Bus
}

// NewRemoveCommandHandler creates a new RemoveCommandHandler
func NewRemoveCommandHandler(r repository.Activity, b event.Bus) *RemoveCommandHandler {
	return &RemoveCommandHandler{
		repo: r,
		bus:  b,
	}
}

func (h RemoveCommandHandler) Invoke(cmd RemoveCommand) error {
	if cmd.ID == "" {
		return exception.NewRequiredField("activity_id")
	}
	return h.persist(cmd.Ctx, cmd.ID)
}

func (h RemoveCommandHandler) persist(ctx context.Context, id string) error {
	snapshot, _, err := h.repo.Fetch(ctx, repository.ActivityCriteria{ID: id})
	if err != nil {
		return err
	} else if err = h.repo.Remove(ctx, snapshot[0].ID()); err != nil {
		return err
	}

	return h.pushEvents(ctx, snapshot[0])
}

func (h RemoveCommandHandler) pushEvents(ctx context.Context, snapshot *aggregate.Activity) error {
	if err := h.bus.Publish(ctx, eventfactory.Activity{}.HardRemoved(snapshot.ID())); err != nil {
		//  rollback
		if errR := h.repo.Save(ctx, *snapshot); errR != nil {
			return errR
		}
		return err
	}

	return nil
}

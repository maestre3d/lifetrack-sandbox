package activity

import (
	"context"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/event"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// AddCommand requests an aggregate.Activity creation
type AddCommand struct {
	Ctx        context.Context
	CategoryID string
	Title      string
}

// AddCommandHandler handles AddCommand(s) requests
type AddCommandHandler struct {
	repo repository.Activity
	bus  event.Bus
}

// NewAddCommandHandler creates a new AddCommandHandler
func NewAddCommandHandler(r repository.Activity, b event.Bus) *AddCommandHandler {
	return &AddCommandHandler{
		repo: r,
		bus:  b,
	}
}

func (h AddCommandHandler) Invoke(cmd AddCommand) (string, error) {
	activity, err := aggregate.NewActivity(cmd.CategoryID, cmd.Title)
	if err != nil {
		return "", err
	} else if err = h.persist(cmd.Ctx, activity); err != nil {
		return "", err
	}

	return activity.ID(), nil
}

func (h AddCommandHandler) persist(ctx context.Context, activity *aggregate.Activity) error {
	if err := h.repo.Save(ctx, *activity); err != nil {
		return err
	}

	return h.pushEvents(ctx, activity)
}

func (h AddCommandHandler) pushEvents(ctx context.Context, activity *aggregate.Activity) error {
	if err := h.bus.Publish(ctx, activity.PullEvents()...); err != nil {
		// rollback
		if errR := h.repo.Remove(ctx, activity.ID()); errR != nil {
			return errR
		}
		return err
	}

	return nil
}

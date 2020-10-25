package category

import (
	"context"
	"strconv"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/event"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// UpdateCommand requests a mutation of an aggregate.Category
type UpdateCommand struct {
	Ctx context.Context
	ID  string

	UserID      string
	Name        string
	Description string
	TargetTime  int64
	Picture     string

	State string
}

// UpdateCommandHandler handles UpdateCommand(s) requests
type UpdateCommandHandler struct {
	repo repository.Category
	bus  event.Bus
}

// NewUpdateCommandHandler creates a new UpdateCommandHandler
func NewUpdateCommandHandler(r repository.Category, b event.Bus) *UpdateCommandHandler {
	return &UpdateCommandHandler{
		repo: r,
		bus:  b,
	}
}

func (h UpdateCommandHandler) Invoke(cmd UpdateCommand) error {
	isCmdEmpty := cmd.Name == "" && cmd.Picture == "" && cmd.State == "" && cmd.UserID == "" && cmd.Description == "" &&
		cmd.TargetTime == 0
	if cmd.ID == "" {
		return exception.NewRequiredField("category_id")
	} else if isCmdEmpty {
		return exception.NewRequiredField("name, description, target_time, picture, state or user_id")
	}

	category, _, err := h.repo.Fetch(cmd.Ctx, repository.CategoryCriteria{ID: cmd.ID})
	if err != nil {
		return err
	}
	snapshot := *category[0]

	if err = h.updater(cmd, category[0]); err != nil {
		return err
	}

	return h.persist(cmd.Ctx, category[0], snapshot)
}

func (h UpdateCommandHandler) updater(cmd UpdateCommand, Category *aggregate.Category) error {
	// update each field
	if cmd.UserID != "" {
		if err := Category.ChangeUser(cmd.UserID); err != nil {
			return err
		}
	}
	if cmd.Name != "" {
		if err := Category.Rename(cmd.Name); err != nil {
			return err
		}
	}
	if cmd.Description != "" {
		if err := Category.ModifyDescription(cmd.Description); err != nil {
			return err
		}
	}
	if cmd.TargetTime > 0 {
		if err := Category.UpdateTargetTime(cmd.TargetTime); err != nil {
			return err
		}
	}
	if cmd.Picture != "" {
		if err := Category.UploadPicture(cmd.Picture); err != nil {
			return err
		}
	}
	if cmd.State != "" {
		if err := h.updateState(cmd.State, Category); err != nil {
			return err
		}
	}

	return nil
}

func (h UpdateCommandHandler) updateState(stateStr string, category *aggregate.Category) error {
	state, err := strconv.ParseBool(stateStr)
	if err != nil {
		return exception.NewFieldFormat("state", "boolean")
	} else if state && !category.State() {
		category.Activate()
	} else if !state && category.State() {
		category.Deactivate()
	}

	return nil
}

func (h UpdateCommandHandler) persist(ctx context.Context, category *aggregate.Category,
	snapshot aggregate.Category) error {
	if err := h.repo.Save(ctx, *category); err != nil {
		return err
	}

	return h.pushEvents(ctx, category, snapshot)
}

func (h UpdateCommandHandler) pushEvents(ctx context.Context, category *aggregate.Category,
	snapshot aggregate.Category) error {
	if err := h.bus.Publish(ctx, category.PullEvents()...); err != nil {
		// rollback
		if errR := h.repo.Save(ctx, snapshot); errR != nil {
			return errR
		}
		return err
	}

	return nil
}

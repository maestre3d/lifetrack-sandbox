package occurrence

import (
	"context"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

type RemoveCommand struct {
	Ctx context.Context
	ID  string
}

type RemoveCommandHandler struct {
	repo repository.Occurrence
}

func NewRemoveCommandHandler(r repository.Occurrence) *RemoveCommandHandler {
	return &RemoveCommandHandler{repo: r}
}

func (h RemoveCommandHandler) Invoke(cmd RemoveCommand) error {
	if cmd.ID == "" {
		return exception.NewRequiredField("occurrence_id")
	}

	// TODO: Publish event to event bus
	return h.repo.Remove(cmd.Ctx, cmd.ID)
}

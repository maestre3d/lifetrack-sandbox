package occurrence

import (
	"context"

	"github.com/alexandria-oss/common-go/exception"
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
}

// NewRemoveCommandHandler creates a RemoveCommandHandler
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

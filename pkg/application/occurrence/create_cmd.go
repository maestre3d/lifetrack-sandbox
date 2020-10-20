package occurrence

import (
	"context"
	"time"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
)

// CreateCommand requests an Occurrence creation
type CreateCommand struct {
	Ctx        context.Context
	ActivityID string
	StartTime  int64
	EndTime    int64
}

// CreateCommandHandler handles CreateCommand(s)
type CreateCommandHandler struct {
	// TODO: Add eventbus
	repo repository.Occurrence
}

// NewCreateCommandHandler
func NewCreateCommandHandler(r repository.Occurrence) *CreateCommandHandler {
	return &CreateCommandHandler{repo: r}
}

// Invoke handle a CreateCommand request, returns Occurrence ID if succeed
func (h CreateCommandHandler) Invoke(cmd CreateCommand) (string, error) {
	occurrence, err := h.creator(cmd.ActivityID, cmd.StartTime, cmd.EndTime)
	if err != nil {
		return "", err
	}

	if err := h.repo.Save(cmd.Ctx, *occurrence); err != nil {
		return "", err
	}
	return occurrence.ID(), nil
}

// creator creates an aggregate.Occurrence parsing required data
func (h CreateCommandHandler) creator(activity string, startTime, endTime int64) (*aggregate.Occurrence, error) {
	return aggregate.NewOccurrence(activity, time.Unix(startTime, 0), time.Unix(endTime, 0))
}

package value

import (
	"time"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/exceptions"
)

const (
	// targetMinDuration target time duration minimum
	targetMinDuration = time.Minute * 5
	// targetMaxDuration target time duration maximum
	targetMaxDuration = time.Hour * 8760
)

// TargetTime expected time duration in a natural year
type TargetTime struct {
	duration time.Duration
}

// NewTargetTime creates a new time duration expected in a natural year
func NewTargetTime(duration time.Duration) (*TargetTime, error) {
	t := new(TargetTime)
	if err := t.SetTarget(duration); err != nil {
		return nil, err
	}
	return t, nil
}

// NewTargetTimeFromPrimitive creates a new time duration expected in a natural year
func NewTargetTimeFromPrimitive(duration time.Duration) *TargetTime {
	return &TargetTime{duration: duration}
}

// SetTarget set the expected duration in a natural year
func (d *TargetTime) SetTarget(duration time.Duration) error {
	memoized := d.duration
	d.duration = duration
	if err := d.IsValid(); err != nil {
		d.duration = memoized
		return err
	}

	return nil
}

// IsValid validates the current Total Duration value
func (d TargetTime) IsValid() error {
	// rule
	//	a.	5 minutes as minimum total duration value; Dom. -> x >= 5 min.
	//	b.	1 year as maximum duration value; Dom. -> x <= 1 year
	if d.duration > targetMaxDuration {
		return exceptions.ErrAboveMaxDuration
	} else if d.duration > time.Nanosecond && d.duration < targetMinDuration {
		return exceptions.ErrBelowMinDuration
	}

	return nil
}

//	--	PRIMITIVES	--

// FromPrimitive sets the current TargetTime value avoiding any domain validation for marshalling purposes
func (d *TargetTime) FromPrimitive(duration time.Duration) {
	d.duration = duration
}

// Duration returns current Duration primitive value
func (d TargetTime) Duration() time.Duration {
	return d.duration
}

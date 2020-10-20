package value

import (
	"time"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/exceptions"
)

const (
	// minDuration http-sandbox total duration minimum
	minDuration = time.Minute * 10
	// maxDuration http-sandbox total duration maximum
	maxDuration = time.Hour * 48
)

// TotalDuration aggregate.Occurrence total time between start time and end time
type TotalDuration struct {
	duration time.Duration
}

// NewTotalDuration creates a new duration between t and u times
func NewTotalDuration(startTime, endTime time.Time) (*TotalDuration, error) {
	t := new(TotalDuration)
	t.Calculate(startTime, endTime)
	if err := t.IsValid(); err != nil {
		return nil, err
	}

	return t, nil
}

// NewTotalDurationFromPrimitive creates a new duration between t and u times
func NewTotalDurationFromPrimitive(duration time.Duration) *TotalDuration {
	return &TotalDuration{duration: duration}
}

// Calculate calculates and sets the duration between start time and end time from the current Occurrence
func (d *TotalDuration) Calculate(startTime, endTime time.Time) {
	d.duration = endTime.Sub(startTime)
}

// IsValid validates the current Total Duration value
func (d TotalDuration) IsValid() error {
	// rule
	//	a.	10 minutes as minimum total duration value; Dom. -> f(x) = start_time - end_time, f(x) >= 10 min.
	//	b.	1 year as maximum duration value; Dom. -> f(x) start_time - end_time, f(x) <= 48 hours
	if d.duration > maxDuration {
		return exceptions.ErrAboveMaxDuration
	} else if d.duration < minDuration {
		return exceptions.ErrBelowMinDuration
	}

	return nil
}

//	--	PRIMITIVES	--

// FromPrimitive sets the current TotalDuration value avoiding any domain validation for marshalling purposes
func (d *TotalDuration) FromPrimitive(duration time.Duration) {
	d.duration = duration
}

// Duration returns current Duration primitive value
func (d TotalDuration) Duration() time.Duration {
	return d.duration
}

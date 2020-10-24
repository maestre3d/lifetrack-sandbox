package aggregate

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/event"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/eventfactory"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/exceptions"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/value"
)

// Occurrence is an event that happened inside each activity, this is a record which helps LifeTrack to keep metrics
// in each activity
type Occurrence struct {
	id            *value.ID
	activity      *value.Activity
	startTime     time.Time
	endTime       time.Time
	totalDuration *value.TotalDuration
	createTime    time.Time
	updateTime    time.Time

	events []event.Domain
}

// NewOccurrence creates a new Occurrence
func NewOccurrence(activityID string, start, end time.Time) (*Occurrence, error) {
	oc := &Occurrence{
		id:            value.NewID(),
		activity:      value.NewActivity(activityID),
		startTime:     start.UTC(),
		endTime:       end.UTC(),
		totalDuration: new(value.TotalDuration),
		createTime:    time.Now().UTC(),
		updateTime:    time.Now().UTC(),
		events:        make([]event.Domain, 0),
	}
	oc.totalDuration.Calculate(oc.startTime, oc.endTime)
	if err := oc.IsValid(); err != nil {
		return nil, err
	}
	oc.RecordEvents(eventfactory.Occurrence{}.ActivityOccurred(*oc.MarshalPrimitive()))

	return oc, nil
}

//	-- BUSINESS CASES	---

// ChangeActivity reassigns occurrence to given activity
func (o *Occurrence) ChangeActivity(activityID string) error {
	memoized := o.activity.String()
	o.activity.FromPrimitive(activityID)
	if err := o.IsValid(); err != nil {
		o.activity.FromPrimitive(memoized)
		return err
	}

	o.updateTime = time.Now().UTC()
	o.RecordEvents(eventfactory.Occurrence{}.Updated(*o.MarshalPrimitive()))
	return nil
}

// EditTimes updates current Occurrence start and end times, then calculates total duration
func (o *Occurrence) EditTimes(startTime, endTime time.Time) error {
	o.startTime = startTime.UTC()
	o.endTime = endTime.UTC()
	o.totalDuration.Calculate(o.startTime, o.endTime)
	if err := o.IsValid(); err != nil {
		return err
	}

	o.updateTime = time.Now().UTC()
	o.RecordEvents(eventfactory.Occurrence{}.Updated(*o.MarshalPrimitive()))
	return nil
}

// IsValid validates the current Occurrence values following domain rules
func (o Occurrence) IsValid() error {
	// rule
	//	a.	5 minutes as minimum total duration value; Dom. -> f(x) = start_time - end_time, f(x) >= 5 min.
	//	b.	48 hours as maximum duration value; Dom. -> f(x) start_time - end_time, f(x) <= 48 hours
	//	c.	1 character as minimum activity value
	if err := o.totalDuration.IsValid(); err != nil {
		return err
	} else if o.activity.IsEmpty() {
		return exceptions.ErrEmptyActivityID
	}

	return nil
}

// --	Event(s)	--

// RecordEvents records given events to the current Occurrence
func (o *Occurrence) RecordEvents(e event.Domain) {
	o.events = append(o.events, e)
}

// PullEvents retrieves all recorded events from the current Occurrence
func (o Occurrence) PullEvents() []event.Domain {
	return o.events
}

// --	MARSHAL / UNMARSHAL	--

// MarshalJSON parses the current Occurrence into a JSON binary
func (o Occurrence) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(*o.MarshalPrimitive())
	if err != nil {
		return nil, exceptions.ErrOccurrenceMarshaling
	}

	return j, nil
}

// MarshalPrimitive parses the current Occurrence into a primitive-only model
func (o Occurrence) MarshalPrimitive() *model.Occurrence {
	return &model.Occurrence{
		ID:            o.id.String(),
		ActivityID:    o.activity.String(),
		StartTime:     o.startTime.Unix(),
		EndTime:       o.endTime.Unix(),
		TotalDuration: int64(o.totalDuration.Duration().Minutes()),
		CreateTime:    o.createTime.Unix(),
		UpdateTime:    o.updateTime.Unix(),
	}
}

// UnmarshalPrimitive creates an Occurrence receiving primitive-only data
//	This function is meant to be used to parse models from databases or any infrastructure storage
func (o *Occurrence) UnmarshalPrimitive(primitive model.Occurrence) error {
	o.id = value.NewIDFromPrimitive(primitive.ID)
	o.activity = value.NewActivity(primitive.ActivityID)
	o.startTime = time.Unix(primitive.StartTime, 0)
	o.endTime = time.Unix(primitive.EndTime, 0)
	o.createTime = time.Unix(primitive.CreateTime, 0)
	o.updateTime = time.Unix(primitive.UpdateTime, 0)

	totalDuration, err := time.ParseDuration(strconv.FormatInt(primitive.TotalDuration, 10) + "m")
	if err != nil {
		return exceptions.ErrInvalidTotalDuration
	}
	o.totalDuration = value.NewTotalDurationFromPrimitive(totalDuration)
	if err := o.IsValid(); err != nil {
		return err
	}

	return nil
}

// Getter(s)

func (o Occurrence) ID() string { return o.id.String() }

func (o Occurrence) Activity() string { return o.activity.String() }

func (o Occurrence) StartTime() time.Time { return o.startTime }

func (o Occurrence) EndTime() time.Time { return o.endTime }

func (o Occurrence) TotalDuration() time.Duration { return o.totalDuration.Duration() }

func (o Occurrence) CreateTime() time.Time { return o.createTime }

func (o Occurrence) UpdateTime() time.Time { return o.updateTime }

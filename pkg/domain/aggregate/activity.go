package aggregate

import (
	"encoding/json"
	"time"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/eventfactory"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"

	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/event"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/exceptions"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/value"
)

// Activity is a task/habit a LifeTrack user wants and will do in further time
type Activity struct {
	id         *value.ID
	category   *value.Category
	title      *value.Title
	picture    *value.Image
	createTime time.Time
	updateTime time.Time
	active     bool

	events []event.Domain
}

// NewActivity creates an Activity with the given values
func NewActivity(categoryID, title string) (*Activity, error) {
	titleV, err := value.NewTitle("", title)
	if err != nil {
		return nil, err
	}

	act := &Activity{
		id:         value.NewID(),
		category:   value.NewCategory(categoryID),
		title:      titleV,
		picture:    value.NewImageFromPrimitive("picture", ""),
		createTime: time.Now().UTC(),
		updateTime: time.Now().UTC(),
		active:     true,
		events:     make([]event.Domain, 0),
	}

	if err = act.IsValid(); err != nil {
		return nil, err
	}
	return act, nil
}

//	--	BUSINESS USE CASES	---

// ChangeCategory changes the current Category
func (a *Activity) ChangeCategory(category string) error {
	memoized := a.category.String()
	a.category.FromPrimitive(category)
	if err := a.IsValid(); err != nil {
		a.category.FromPrimitive(memoized)
		return err
	}
	a.updateTime = time.Now().UTC()
	a.RecordEvents(eventfactory.Activity{}.Updated(*a.MarshalPrimitive()))
	return nil
}

// Rename changes the current Activity title
func (a *Activity) Rename(title string) error {
	if err := a.title.Rename(title); err != nil {
		return err
	} else if err := a.IsValid(); err != nil {
		return err
	}
	a.updateTime = time.Now().UTC()
	a.RecordEvents(eventfactory.Activity{}.Updated(*a.MarshalPrimitive()))
	return nil
}

// UploadPicture changes the current Activity picture
func (a *Activity) UploadPicture(picture string) error {
	if err := a.picture.Save(picture); err != nil {
		return err
	}
	a.updateTime = time.Now().UTC()
	a.RecordEvents(eventfactory.Activity{}.Updated(*a.MarshalPrimitive()))
	return nil
}

// Activate sets the current Activity state to active
func (a *Activity) Activate() {
	a.active = true
	a.updateTime = time.Now().UTC()
	a.RecordEvents(eventfactory.Activity{}.Restored(a.id.String()))
}

// Deactivate sets the current Activity state to deactivated
func (a *Activity) Deactivate() {
	a.active = false
	a.updateTime = time.Now().UTC()
	a.RecordEvents(eventfactory.Activity{}.Removed(a.id.String()))
}

// IsValid validates the current Activity state
func (a Activity) IsValid() error {
	//	rules
	//	a.	title min length 2 character
	//	b.	category is required
	if a.category.IsEmpty() {
		return exceptions.ErrEmptyCategoryID
	} else if len(a.title.String()) < 2 {
		return exceptions.ErrBelowTitleLength
	}

	return nil
}

// --	Event(s)	--

// RecordEvents records given events to the current Activity
func (a *Activity) RecordEvents(e event.Domain) {
	a.events = append(a.events, e)
}

// PullEvents retrieves all recorded events from the current Activity
func (a Activity) PullEvents() []event.Domain {
	return a.events
}

// --	MARSHAL / UNMARSHAL	--

// MarshalJSON parses the current Activity into a JSON binary
func (a Activity) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(*a.MarshalPrimitive())
	if err != nil {
		return nil, exceptions.ErrActivityMarshaling
	}

	return j, nil
}

// MarshalPrimitive parses the current Activity into a primitive-only model
func (a Activity) MarshalPrimitive() *model.Activity {
	return &model.Activity{
		ID:         a.id.String(),
		CategoryID: a.category.String(),
		Title:      a.title.String(),
		Picture:    a.picture.String(),
		CreateTime: a.createTime.Unix(),
		UpdateTime: a.updateTime.Unix(),
		Active:     a.active,
	}
}

// UnmarshalPrimitive creates an Activity receiving primitive-only data
//	This function is meant to be used to parse models from databases or any infrastructure storage
func (a *Activity) UnmarshalPrimitive(primitive model.Activity) error {
	a.id = value.NewIDFromPrimitive(primitive.ID)
	a.category = value.NewCategory(primitive.CategoryID)
	a.title = value.NewTitleFromPrimitive("", primitive.Title)
	a.picture = value.NewImageFromPrimitive("picture", primitive.Picture)
	a.createTime = time.Unix(primitive.CreateTime, 0)
	a.updateTime = time.Unix(primitive.UpdateTime, 0)
	a.active = primitive.Active
	if err := a.IsValid(); err != nil {
		return err
	}

	return nil
}

//	--	GETTER(S)	--

func (a Activity) ID() string { return a.id.String() }

func (a Activity) Category() string { return a.category.String() }

func (a Activity) Title() string { return a.title.String() }

func (a Activity) Picture() string { return a.picture.String() }

func (a Activity) CreateTime() time.Time { return a.createTime }

func (a Activity) UpdateTime() time.Time { return a.updateTime }

func (a Activity) State() bool { return a.active }

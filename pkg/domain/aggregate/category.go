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

type Category struct {
	id   *value.ID
	user *value.User

	name        *value.Title
	description *value.Description
	targetTime  *value.TargetTime
	picture     *value.Image

	createTime time.Time
	updateTime time.Time
	active     bool

	events []event.Domain
}

// NewCategory creates a Category with the given values
func NewCategory(userID, name string) (*Category, error) {
	nameV, err := value.NewTitle("name", name)
	if err != nil {
		return nil, err
	}

	cat := &Category{
		id:          value.NewID(),
		user:        value.NewUser(userID),
		name:        nameV,
		description: &value.Description{},
		targetTime:  &value.TargetTime{},
		picture:     value.NewImageFromPrimitive("picture", ""),
		createTime:  time.Now().UTC(),
		updateTime:  time.Now().UTC(),
		active:      true,
		events:      make([]event.Domain, 0),
	}

	if err = cat.IsValid(); err != nil {
		return nil, err
	}
	return cat, nil
}

//	--	BUSINESS USE CASES	---

// ChangeUser changes the current User
func (c *Category) ChangeUser(user string) error {
	memoized := c.user.String()
	c.user.FromPrimitive(user)
	if err := c.IsValid(); err != nil {
		c.user.FromPrimitive(memoized)
		return err
	}
	c.updateTime = time.Now().UTC()
	c.RecordEvents(eventfactory.Category{}.Updated(*c.MarshalPrimitive()))
	return nil
}

// Rename changes the current Category name
func (c *Category) Rename(name string) error {
	if err := c.name.Rename(name); err != nil {
		return err
	} else if err := c.IsValid(); err != nil {
		return err
	}
	c.updateTime = time.Now().UTC()
	c.RecordEvents(eventfactory.Category{}.Updated(*c.MarshalPrimitive()))
	return nil
}

// ModifyDescription update the current Category description
func (c *Category) ModifyDescription(description string) error {
	if err := c.description.Change(description); err != nil {
		return err
	} else if err := c.IsValid(); err != nil {
		return err
	}
	c.updateTime = time.Now().UTC()
	c.RecordEvents(eventfactory.Category{}.Updated(*c.MarshalPrimitive()))
	return nil
}

// UpdateTargetTime set a new target time to the current Category
func (c *Category) UpdateTargetTime(targetTime int64) error {
	targetTimeV, err := time.ParseDuration(strconv.FormatInt(targetTime, 10) + "m")
	if err != nil {
		return exceptions.ErrInvalidTargetTime
	} else if err = c.targetTime.SetTarget(targetTimeV); err != nil {
		return err
	}
	c.updateTime = time.Now().UTC()
	c.RecordEvents(eventfactory.Category{}.Updated(*c.MarshalPrimitive()))
	return nil
}

// UploadPicture changes the current Category picture
func (c *Category) UploadPicture(picture string) error {
	if err := c.picture.Save(picture); err != nil {
		return err
	}
	c.updateTime = time.Now().UTC()
	c.RecordEvents(eventfactory.Category{}.Updated(*c.MarshalPrimitive()))
	return nil
}

// Activate sets the current Category state to active
func (c *Category) Activate() {
	c.active = true
	c.updateTime = time.Now().UTC()
	c.RecordEvents(eventfactory.Category{}.Restored(c.id.String()))
}

// Deactivate sets the current Category state to deactivated
func (c *Category) Deactivate() {
	c.active = false
	c.updateTime = time.Now().UTC()
	c.RecordEvents(eventfactory.Category{}.Removed(c.id.String()))
}

// IsValid validates the current Category state
func (c Category) IsValid() error {
	//	rules
	//	a.	name min length 2 character
	//	b.	user is required
	if c.user.IsEmpty() {
		return exceptions.ErrEmptyCategoryID
	} else if len(c.name.String()) < 2 {
		return exceptions.ErrBelowNameLength
	}

	return nil
}

// --	Event(s)	--

// RecordEvents records given events to the current Category
func (c *Category) RecordEvents(e event.Domain) {
	c.events = append(c.events, e)
}

// PullEvents retrieves all recorded events from the current Category
func (c Category) PullEvents() []event.Domain {
	return c.events
}

// --	MARSHAL / UNMARSHAL	--

// MarshalJSON parses the current Category into c JSON binary
func (c Category) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(*c.MarshalPrimitive())
	if err != nil {
		return nil, exceptions.ErrCategoryMarshaling
	}

	return j, nil
}

// MarshalPrimitive parses the current Category into c primitive-only model
func (c Category) MarshalPrimitive() *model.Category {
	return &model.Category{
		ID:          c.id.String(),
		UserID:      c.user.String(),
		Name:        c.name.String(),
		Description: c.description.String(),
		TargetTime:  int64(c.targetTime.Duration().Minutes()),
		Picture:     c.picture.String(),
		CreateTime:  c.createTime.Unix(),
		UpdateTime:  c.updateTime.Unix(),
		Active:      c.active,
	}
}

// UnmarshalPrimitive creates an Category receiving primitive-only data
//	This function is meant to be used to parse models from databases or any infrastructure storage
func (c *Category) UnmarshalPrimitive(primitive model.Category) error {
	c.id = value.NewIDFromPrimitive(primitive.ID)
	c.user = value.NewUser(primitive.UserID)
	c.name = value.NewTitleFromPrimitive("name", primitive.Name)
	c.description = value.NewDescriptionFromPrimitive("", primitive.Description)
	c.picture = value.NewImageFromPrimitive("picture", primitive.Picture)
	c.createTime = time.Unix(primitive.CreateTime, 0)
	c.updateTime = time.Unix(primitive.UpdateTime, 0)
	c.active = primitive.Active

	targetTime, err := time.ParseDuration(strconv.FormatInt(primitive.TargetTime, 10) + "m")
	if err != nil {
		return exceptions.ErrInvalidTargetTime
	}
	c.targetTime = value.NewTargetTimeFromPrimitive(targetTime)

	if err := c.IsValid(); err != nil {
		return err
	}

	return nil
}

//	--	GETTER(S)	--

func (c Category) ID() string { return c.id.String() }

func (c Category) User() string { return c.user.String() }

func (c Category) Name() string { return c.name.String() }

func (c Category) Description() string { return c.description.String() }

func (c Category) TargetTime() time.Duration { return c.targetTime.Duration() }

func (c Category) Picture() string { return c.picture.String() }

func (c Category) CreateTime() time.Time { return c.createTime }

func (c Category) UpdateTime() time.Time { return c.updateTime }

func (c Category) State() bool { return c.active }

package dynamocategory

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/model"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/persistence/dynamoutils"
)

// categoryDynamo custom AWS DynamoDB model for model.Category/aggregate.Category
type categoryDynamo struct {
	// Category ID
	PK string `json:"PK"`
	// Category ID
	SK string `json:"SK"`
	// User ID
	GSIPK string `json:"GSIPK"`
	// Category ID
	GSISK       string `json:"GSISK"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	TargetTime  int64  `json:"target_time,omitempty"`
	Picture     string `json:"picture,omitempty"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	Active      bool   `json:"active"`
}

//	--	GETTERS	--

// ID get the current categoryDynamo unique identifier
func (m categoryDynamo) ID() string {
	return dynamoutils.DecomposeCompositeKey(m.PK)
}

// User get the current categoryDynamo user reference
func (m categoryDynamo) User() string {
	return dynamoutils.DecomposeCompositeKey(m.GSIPK)
}

//	-- MARSHAL/UN-MARSHAL	---

// MarshalModel retrieves the current categoryDynamo as a model.Category
func (m categoryDynamo) MarshalModel() *model.Category {
	return &model.Category{
		ID:          m.ID(),
		UserID:      m.User(),
		Name:        m.Name,
		Description: m.Description,
		TargetTime:  m.TargetTime,
		Picture:     m.Picture,
		CreateTime:  m.CreateTime,
		UpdateTime:  m.UpdateTime,
		Active:      m.Active,
	}
}

// MarshalModel retrieves the current categoryDynamo as an aggregate.Category
func (m categoryDynamo) MarshalAggregate() (*aggregate.Category, error) {
	c := new(aggregate.Category)
	err := c.UnmarshalPrimitive(*m.MarshalModel())
	if err != nil {
		return nil, err
	}

	return c, nil
}

// UnmarshalAggregate parses the given aggregate.Category into the current categoryDynamo
func (m *categoryDynamo) UnmarshalAggregate(category aggregate.Category) {
	primitive := category.MarshalPrimitive()
	compositeKey := dynamoutils.NewCompositeKey(schemaName, primitive.ID)

	m.PK = compositeKey
	m.SK = compositeKey
	m.GSIPK = dynamoutils.NewCompositeKey("user", primitive.UserID)
	m.GSISK = compositeKey
	m.Name = primitive.Name
	m.CreateTime = primitive.CreateTime
	m.UpdateTime = primitive.UpdateTime
	m.Active = primitive.Active

	// Optional
	m.Description = primitive.Description
	m.TargetTime = primitive.TargetTime
	m.Picture = primitive.Picture
}

// MarshalAttributeValues parses the current categoryDynamo values into a dynamo attribute value map
func (m categoryDynamo) MarshalAttributeValues() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		":n": {
			S: aws.String(m.Name),
		},
		":d": {
			S: aws.String(m.Description),
		},
		":tt": {
			N: aws.String(strconv.FormatInt(m.TargetTime, 10)),
		},
		":p": {
			S: aws.String(m.Picture),
		},
		":ct": {
			N: aws.String(strconv.FormatInt(m.CreateTime, 10)),
		},
		":ut": {
			N: aws.String(strconv.FormatInt(m.UpdateTime, 10)),
		},
		":a": {
			BOOL: aws.Bool(m.Active),
		},
		":u": {
			S: aws.String(m.GSIPK),
		},
		":s": {
			S: aws.String(m.GSISK),
		},
	}
}

func (m categoryDynamo) MarshalKeys() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"PK": {
			S: aws.String(m.PK),
		},
		"SK": {
			S: aws.String(m.SK),
		},
	}
}

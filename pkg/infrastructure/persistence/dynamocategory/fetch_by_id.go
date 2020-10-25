package dynamocategory

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/configuration"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/persistence/dynamoutils"
)

// fetchByID strategy when criteria contains a category ID (ID)
type fetchByID struct {
	cfg configuration.Configuration
	db  *dynamodb.DynamoDB
}

func (r fetchByID) Do(ctx context.Context, criteria repository.CategoryCriteria) ([]*aggregate.Category, string, error) {
	compositeKey := aws.String(dynamoutils.NewCompositeKey(schemaName, criteria.ID))
	o, err := r.db.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: compositeKey,
			},
			"SK": {
				S: compositeKey,
			},
		},
		TableName: aws.String(r.cfg.DynamoTable.Name),
	})
	if err != nil {
		return nil, "", getDomainError(err, r.cfg)
	}

	m := new(categoryDynamo)
	err = dynamodbattribute.UnmarshalMap(o.Item, m)
	if err != nil {
		return nil, "", getDomainError(err, r.cfg)
	}
	c, err := m.MarshalAggregate()
	if err != nil {
		return nil, "", err
	}

	return []*aggregate.Category{c}, "", nil
}

package dynamocategory

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/aggregate"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/exceptions"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/repository"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/configuration"
)

const schemaName = "category"

var (
	dynamoSingleton *Dynamo
	dynamoOnce      = new(sync.Once)
)

type Dynamo struct {
	cfg configuration.Configuration
	db  *dynamodb.DynamoDB
	mu  *sync.RWMutex
}

func NewDynamo(cfg configuration.Configuration, pool *dynamodb.DynamoDB) *Dynamo {
	dynamoOnce.Do(func() {
		dynamoSingleton = &Dynamo{
			cfg: cfg,
			db:  pool,
			mu:  new(sync.RWMutex),
		}
	})

	return dynamoSingleton
}

func (r *Dynamo) Save(ctx context.Context, category aggregate.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	m := new(categoryDynamo)
	m.UnmarshalAggregate(category)

	_, err := r.db.UpdateItemWithContext(ctx, &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#N":  aws.String("name"),
			"#D":  aws.String("description"),
			"#TT": aws.String("target_time"),
			"#P":  aws.String("picture"),
			"#CT": aws.String("create_time"),
			"#UT": aws.String("update_time"),
			"#A":  aws.String("active"),
			"#U":  aws.String("GSIPK"),
			"#S":  aws.String("GSISK"),
		},
		ExpressionAttributeValues: m.MarshalAttributeValues(),
		Key:                       m.MarshalKeys(),
		ReturnValues:              aws.String(dynamodb.ReturnValueNone),
		TableName:                 aws.String(r.cfg.DynamoTable.Name),
		UpdateExpression: aws.String("SET #N = :n, #D = :d, #TT = :tt, #P = :p, #CT = :ct, #UT = :ut, " +
			"#A = :a, #U = :u, #S = :s"),
	})

	return getDomainError(err, r.cfg)
}

func (r *Dynamo) Fetch(ctx context.Context, criteria repository.CategoryCriteria) ([]*aggregate.Category, string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fetch := r.setFetchStrategy(criteria)
	cats, nextToken, err := fetch.Do(ctx, criteria)
	if len(cats) == 0 && err == nil {
		return nil, "", exceptions.ErrCategoryNotFound
	}

	return cats, nextToken, err
}

func (r *Dynamo) Remove(_ context.Context, _ string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	panic("not implemented")
}

func (r Dynamo) setFetchStrategy(criteria repository.CategoryCriteria) fetchStrategy {
	switch {
	case criteria.ID != "":
		return fetchByID{cfg: r.cfg, db: r.db}
	default:
		return nil
	}
}

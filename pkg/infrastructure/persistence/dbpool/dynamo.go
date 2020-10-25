package dbpool

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/configuration"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/remote"
)

// NewDynamoDBPool creates a new AWS DynamoDB connection pool
func NewDynamoDBPool(cfg configuration.Configuration) *dynamodb.DynamoDB {
	return dynamodb.New(remote.NewAWSSession(cfg.DynamoTable.Region))
}

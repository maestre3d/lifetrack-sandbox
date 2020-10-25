package dynamocategory

import (
	"github.com/alexandria-oss/common-go/exception"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/maestre3d/lifetrack-sanbox/pkg/domain/exceptions"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/configuration"
)

// getDomainError returns a valid domain error from awserr.Error dynamodb exceptions
func getDomainError(err error, cfg configuration.Configuration) error {
	if err != nil {
		if errAWS, ok := err.(awserr.Error); ok {
			switch errAWS.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return exceptions.ErrCategoryNotFound
			case dynamodb.ErrCodeIndexNotFoundException:
				return exception.NewNotFound("category_id")
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return exception.NewFieldFormat("category_conditional", "valid query conditional field")
			case dynamodb.ErrCodeRequestLimitExceeded:
				return exception.NewNetworkCall("aws dynamodb table "+cfg.DynamoTable.Name,
					cfg.DynamoTable.Region)
			}
		}

		return err
	}

	return nil
}

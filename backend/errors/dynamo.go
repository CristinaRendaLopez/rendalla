package errors

import (
	stdErrors "errors"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
)

// Constants for AWS error codes not included in the DynamoDB SDK constants.
const (
	AWSValidationException             = "ValidationException"
	AWSAccessDeniedException           = "AccessDeniedException"
	AWSThrottlingException             = "ThrottlingException"
	AWSRequestLimitExceeded            = "RequestLimitExceeded"
	AWSTransactionConflictException    = "TransactionConflictException"
	AWSItemCollectionSizeLimitExceeded = "ItemCollectionSizeLimitExceededException"
)

// HandleDynamoError maps low-level DynamoDB and AWS SDK errors to custom application errors.
// Returns a more meaningful application-level error based on the source of failure.
func HandleDynamoError(err error) error {
	if err == nil {
		return nil
	}

	if stdErrors.Is(err, dynamo.ErrNotFound) {
		return ErrResourceNotFound
	}

	if stdErrors.Is(err, dynamo.ErrTooMany) {
		return ErrTooManyResults
	}

	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case dynamodb.ErrCodeConditionalCheckFailedException:
			return ErrOperationNotAllowed
		case dynamodb.ErrCodeProvisionedThroughputExceededException,
			AWSThrottlingException,
			AWSRequestLimitExceeded:
			return ErrThroughputExceeded
		case dynamodb.ErrCodeResourceNotFoundException:
			return ErrResourceNotFound
		case dynamodb.ErrCodeInternalServerError:
			return ErrInternalServer
		case AWSAccessDeniedException:
			return ErrUnauthorized
		case AWSValidationException:
			return ErrBadRequest
		case AWSItemCollectionSizeLimitExceeded:
			return ErrValidationFailed
		case AWSTransactionConflictException:
			return ErrOperationNotAllowed
		default:
			return ErrInternalServer
		}
	}

	return ErrInternalServer
}

package repository

// PagingKey represents an abstract pagination token used to continue a paginated query.
// The concrete type depends on the underlying database implementation.
// For example, in DynamoDB it typically wraps a LastEvaluatedKey.
type PagingKey interface{}

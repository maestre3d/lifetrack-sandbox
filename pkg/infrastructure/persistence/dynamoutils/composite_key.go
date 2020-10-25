package dynamoutils

import (
	"fmt"
	"strings"
)

// NewCompositeKey returns a formatted Composite Key for AWS DynamoDB tables
func NewCompositeKey(schema, id string) string {
	return fmt.Sprintf("%s#%s", strings.ToLower(schema), id)
}

// DecomposeCompositeKey returns a decomposed composite key, if non-valid returns the given key
func DecomposeCompositeKey(key string) string {
	if !strings.Contains(key, "#") {
		return key
	}

	return strings.Split(key, "#")[1]
}

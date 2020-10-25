package configuration

// Configuration kernel configuration
type Configuration struct {
	Version     string      `json:"version"`
	Stage       string      `json:"stage"`
	DynamoTable dynamoTable `json:"dynamo_table"`
}

// dynamoTable AWS DynamoDB Table config
type dynamoTable struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

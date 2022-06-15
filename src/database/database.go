// Package database handles transactions made to DynamoDB
package database

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
	"strings"
)

const (
	tableNameEnvVar = "GITHUB_TABLE_NAME"
)

var dynamoDbClient *dynamodb.DynamoDB

type scanCallback = func(item map[string]*dynamodb.AttributeValue) error

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamoDbClient = dynamodb.New(sess)
}

func TableName() *string {
	return aws.String(os.Getenv(tableNameEnvVar))
}

// DBClient returns the connection client to DynamoDB
func DBClient() *dynamodb.DynamoDB {
	return dynamoDbClient
}

// ScanFilterByType scans the table and executes operation for each item. The `cb` or callback provided will be executed for items with Type property that
// that matches the type provided in `filterByType`
func ScanFilterByType(filterByType string, cb scanCallback) error {
	input := &dynamodb.ScanInput{
		TableName: TableName(),
	}
	var lastEvaluated map[string]*dynamodb.AttributeValue

	for {
		if lastEvaluated != nil {
			input.ExclusiveStartKey = lastEvaluated
		}

		output, err := DBClient().Scan(input)
		fmt.Printf("Count: %d \n", aws.Int64Value(output.Count))
		fmt.Printf("ScannedCount: %d \n", aws.Int64Value(output.ScannedCount))

		if err != nil {
			return err
		}

		for _, item := range output.Items {
			if filterByType != "" && aws.StringValue(item["Type"].S) != filterByType {
				continue
			}

			if err := cb(item); err != nil {
				return err
			}
		}

		if output.LastEvaluatedKey != nil {
			lastEvaluated = output.LastEvaluatedKey
		} else {
			break
		}

	}

	return nil
}

func GenerateAttrNotExistsExpression(fields ...string) *string {
	sb := strings.Builder{}

	for index, field := range fields {
		if index != 0 {
			sb.WriteString(" AND ")
		}

		stm := fmt.Sprintf("attribute_not_exists(%s)", field)
		sb.WriteString(stm)
	}

	return aws.String(sb.String())
}

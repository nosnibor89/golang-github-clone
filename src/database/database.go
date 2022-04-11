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

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamoDbClient = dynamodb.New(sess)
}

func tableName() *string {
	return aws.String(os.Getenv(tableNameEnvVar))
}

func generateAttrNotExistsExpression(fields ...string) *string {
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

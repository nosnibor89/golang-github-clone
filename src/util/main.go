package util

import (
	"fmt"
	"github-clone/src/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"strings"
)

type HttpError struct {
	Code    int
	Message string
}

// MergeMaps Merge multiple maps into one. Key values will be overridden in the same order the maps were provided
func MergeMaps[K string | int64, V any](maps ...map[K]V) map[K]V {
	newMap := make(map[K]V)

	for _, currentMap := range maps {
		for key, value := range currentMap {
			newMap[key] = value
		}
	}

	return newMap
}

func StringIsEmpty(text string) bool {
	return len(text) == 0
}

func GetUserFromRequest(request events.APIGatewayProxyRequest) model.User {
	user := request.Headers["Authorization"]
	return model.User{
		Username: user,
		Name:     user,
	}
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

func HttpErrorFromException(e error) HttpError {
	if strings.Contains(e.Error(), "ConditionalCheckFailedException:") {
		return HttpError{
			Code:    400,
			Message: "Entity already exists",
		}
	}
	return HttpError{
		Code:    500,
		Message: "Error creating entity",
	}
}

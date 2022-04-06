package util

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"strings"
	"time"
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

func ParseTimeAttr(datetime string) time.Time {
	parsed, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		fmt.Printf("Count not parse time for value %v\n", datetime)
		return time.Date(1997, time.January, 1, 1, 1, 1, 1, time.UTC)
	}

	return parsed
}

func ParseTimeItem(datetime time.Time) string {
	return datetime.Format(time.RFC3339)
}

//TODO: Move this to it's own package/file
var reset = "\033[0m"
var red = "\033[31m"

//var green = "\033[32m"
//var yellow = "\033[33m"
//var blue = "\033[34m"
//var purple = "\033[35m"
var cyan = "\033[36m"

//var gray = "\033[37m"
//var white = "\033[97m"

func PrintRed(data interface{}) {
	fmt.Printf("%s %v %s\n", red, data, reset)
}

func PrintCyan(data interface{}) {
	fmt.Printf("%s %v %s\n", cyan, data, reset)
}

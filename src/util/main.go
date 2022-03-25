package util

import (
	"github-clone/src/model"
	"github.com/aws/aws-lambda-go/events"
)

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

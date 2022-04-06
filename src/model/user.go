package model

import "github.com/aws/aws-lambda-go/events"

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

func GetUserFromRequest(request events.APIGatewayProxyRequest) User {
	user := request.Headers["Authorization"]
	return User{
		Username: user,
		Name:     user,
	}
}

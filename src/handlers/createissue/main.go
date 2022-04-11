package main

import (
	"context"
	"fmt"
	"github-clone/src/database"
	"github-clone/src/errors"
	"github-clone/src/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body string
	statusCode := 201

	user := model.GetUserFromRequest(request)

	issue := model.Issue{}

	if err := issue.FromJSON(request.Body); err != nil {
		httpError := errors.HttpErrorFromException(err)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	issue.Repo = model.Repo{
		Name: request.PathParameters["repo"],
		Owner: model.User{
			Username: request.PathParameters["owner"],
			Name:     request.PathParameters["owner"],
		},
	}

	issue.WithCreator(user)

	repository := database.Issue{}
	newIssue, err := repository.Create(issue)

	if err != nil {
		httpError := errors.HttpErrorFromException(err)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	body, decodingError := newIssue.ToJSON()

	if decodingError == nil {
		statusCode = 200
		fmt.Println("A new repo was created")
	}

	return events.APIGatewayProxyResponse{Body: body, StatusCode: statusCode}, nil
}

func main() {
	lambda.Start(handleRequest)
}

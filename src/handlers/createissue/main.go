package main

import (
	"context"
	dbIssue "github-clone/src/database/issue"
	"github-clone/src/errors"
	"github-clone/src/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

var (
	issue = model.Issue{}
)

func handleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body string
	statusCode := http.StatusCreated

	user := model.GetUserFromRequest(request)

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

	newIssue, err := dbIssue.Create(issue)

	if err != nil {
		httpError := errors.HttpErrorFromException(err)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	body, decodingError := newIssue.ToJSON()

	if decodingError != nil {
		httpError := errors.HttpErrorFromException(decodingError)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	return events.APIGatewayProxyResponse{Body: body, StatusCode: statusCode}, nil
}

func main() {
	lambda.Start(handleRequest)
}

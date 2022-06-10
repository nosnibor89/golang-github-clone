package main

import (
	"context"
	"github-clone/src/database/pullrequest"
	"github-clone/src/errors"
	"github-clone/src/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

var (
	pullRequest = model.PullRequest{}
)

func handleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body string
	statusCode := http.StatusCreated

	user := model.GetUserFromRequest(request)

	if err := pullRequest.FromJSON(request.Body); err != nil {
		httpError := errors.HttpErrorFromException(err)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	pullRequest.Repo = model.Repo{
		Name: request.PathParameters["repo"],
		Owner: model.User{
			Username: request.PathParameters["owner"],
			Name:     request.PathParameters["owner"],
		},
	}

	pullRequest.WithCreator(user)

	newPullRequest, err := pullrequest.Create(pullRequest)

	if err != nil {
		httpError := errors.HttpErrorFromException(err)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	body, decodingError := newPullRequest.ToJSON()

	if decodingError != nil {
		httpError := errors.HttpErrorFromException(decodingError)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	return events.APIGatewayProxyResponse{Body: body, StatusCode: statusCode}, nil
}

func main() {
	lambda.Start(handleRequest)
}

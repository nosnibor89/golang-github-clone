package main

import (
	"context"
	"github-clone/src/database/stars"
	"github-clone/src/errors"
	"github-clone/src/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

func handleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repo := request.PathParameters["repo"]
	owner := request.PathParameters["owner"]

	user := model.GetUserFromRequest(request)

	if err := stars.UnStarRepo(repo, owner, user.Username); err != nil {
		httpError := errors.HttpErrorFromException(err)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}

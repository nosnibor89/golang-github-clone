package main

import (
	"context"
	"encoding/json"
	"github-clone/src/database/pullrequest"
	"github-clone/src/errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

func handleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repo := request.PathParameters["repo"]
	owner := request.PathParameters["owner"]
	status := request.QueryStringParameters["status"]

	foundRepo, issues, err := pullrequest.Find(repo, owner, status)

	if err != nil {
		httpError := errors.HttpErrorFromException(err)
		return events.APIGatewayProxyResponse{
			StatusCode: httpError.Code,
			Body:       httpError.Message,
		}, nil
	}

	body, _ := json.Marshal(map[string]interface{}{
		"repo":         foundRepo,
		"pullRequests": issues,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}

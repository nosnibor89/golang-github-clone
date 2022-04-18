package main

import (
	"context"
	"encoding/json"
	"github-clone/src/database"
	"github-clone/src/errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

var issue = database.Issue{}

func handleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repo := request.PathParameters["repo"]
	owner := request.PathParameters["owner"]
	status := request.QueryStringParameters["status"]

	foundRepo, issues, err := issue.GetIssues(repo, owner, status)

	if err != nil {
		httpError := errors.HttpErrorFromException(err)
		return events.APIGatewayProxyResponse{
			StatusCode: httpError.Code,
			Body:       httpError.Message,
		}, nil
	}

	body, _ := json.Marshal(map[string]interface{}{
		"repo":   foundRepo,
		"issues": issues,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}

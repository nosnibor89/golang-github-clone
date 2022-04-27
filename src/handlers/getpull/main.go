package main

import (
	"context"
	"fmt"
	"github-clone/src/database"
	"github-clone/src/errors"
	"github-clone/src/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

var (
	pullRequest = model.PullRequest{}
	db          = database.PullRequest{}
)

func handleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body string
	statusCode := http.StatusOK

	repo := request.PathParameters["repo"]
	prNumber := request.PathParameters["prNumber"]

	findOneInput := database.PullRequestFindOneInput{
		Repo:              repo,
		Owner:             request.PathParameters["owner"],
		PullRequestNumber: prNumber,
	}

	found, err := db.FindOne(findOneInput)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       fmt.Sprintf("Could not find pull request %s from repo %s", prNumber, repo),
		}, nil
	}

	body, decodingError := found.ToJSON()

	if decodingError != nil {
		httpError := errors.HttpErrorFromException(decodingError)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	return events.APIGatewayProxyResponse{Body: body, StatusCode: statusCode}, nil
}

func main() {
	lambda.Start(handleRequest)
}

package main

import (
	"context"
	"fmt"
	"github-clone/src/database/issue"
	"github-clone/src/errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
)

/*
	At this point any user can ask for any other user's repo which is what GitHub does, but maybe we need to handle that in a different way
*/
func handleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repo := request.PathParameters["repo"]
	owner := request.PathParameters["owner"]
	issueNumber := request.PathParameters["issueNumber"]

	input := issue.FindIssueInput{
		Repo:        repo,
		Owner:       owner,
		IssueNumber: issueNumber,
	}

	found, err := issue.FindIssue(input)

	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       fmt.Sprintf("Could not find issue %s from repo %s", issueNumber, repo),
		}, nil
	}

	body, decodingError := found.ToJSON()

	if decodingError != nil {
		httpError := errors.HttpErrorFromException(decodingError)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       body,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}

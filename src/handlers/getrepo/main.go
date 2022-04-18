package main

import (
	"context"
	"github-clone/src/database"
	"github-clone/src/errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
)

var db = database.Repository{}

/*
	At this point any user can ask for any other user's repo which is what GitHub does, but maybe we need to handle that in a different way
*/
func handleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repo := request.PathParameters["repo"]
	owner := request.PathParameters["owner"]

	found := db.FindOne(repo, owner)

	if found == nil {
		log.Printf("repo with owner %s and name %s was not found", owner, repo)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "Could not find github repo",
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

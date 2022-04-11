package main

import (
	"context"
	"encoding/json"
	"github-clone/src/database"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

//TODO: At this point any use can ask for any other user's repo which is what github does but maybe we need to handle that in a different way

func handleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	repository := database.Repository{}

	repo := request.PathParameters["repo"]
	owner := request.PathParameters["owner"]

	found := repository.FindOne(repo, owner)

	if found == nil {
		log.Printf("repo with owner %s and name %s was not found", owner, repo)
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Could not find github repo",
		}, nil
	}

	body, _ := json.Marshal(found)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}

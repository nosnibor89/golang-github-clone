package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github-clone/src/model"
	"github-clone/src/repositories"
	"github-clone/src/util"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	//fmt.Printf("Body size = %d.\n", len(request.Body))
	//fmt.Println("Headers:")
	//for key, value := range request.Headers {
	//	fmt.Printf("    %s: %s\n", key, value)
	//}
	var requestError error
	var body string
	statusCode := 201

	user := util.GetUserFromRequest(request)

	repo := model.Repo{}

	if err := json.Unmarshal([]byte(request.Body), &repo); err != nil {
		msg := fmt.Sprintf("Could not parse body correctly %v\n", err)
		fmt.Println(msg)
		requestError = err
		statusCode = 400
		body = msg
	}

	repo.Owner = user

	repository := repositories.GithubRepository{}
	newRepo, err := repository.Create(repo)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Error creating repo", StatusCode: 500}, err
	}

	if decoded, err := json.Marshal(newRepo); err != nil {
		msg := fmt.Sprintf("Could not parse created repo %v\n", err)
		fmt.Println(msg)
		requestError = err
		statusCode = 500
		body = msg
	} else {
		fmt.Println("A new repo was created")
		body = string(decoded)
	}

	return events.APIGatewayProxyResponse{Body: body, StatusCode: statusCode}, requestError
}

func main() {
	lambda.Start(handleRequest)
}

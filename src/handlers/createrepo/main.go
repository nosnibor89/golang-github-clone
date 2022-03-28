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
	var body string
	statusCode := 201

	user := util.GetUserFromRequest(request)

	repo := model.Repo{}

	//TODO: Find re-usability for error handling
	if err := json.Unmarshal([]byte(request.Body), &repo); err != nil {
		msg := fmt.Sprintf("Could not parse body correctly %v\n", err)
		fmt.Println(msg)
		statusCode = util.HttpErrorFromException(err).Code
		statusCode = 400
		body = msg
	}

	repo.Owner = user

	repository := repositories.GithubRepository{}
	newRepo, err := repository.Create(repo)

	if err != nil {
		httpError := util.HttpErrorFromException(err)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	if decoded, err := json.Marshal(newRepo); err != nil {
		msg := fmt.Sprintf("Could not parse created repo %v\n", err)
		fmt.Println(msg)
		statusCode = util.HttpErrorFromException(err).Code
		body = msg
	} else {
		fmt.Println("A new repo was created")
		body = string(decoded)
	}

	return events.APIGatewayProxyResponse{Body: body, StatusCode: statusCode}, nil
}

func main() {
	lambda.Start(handleRequest)
}

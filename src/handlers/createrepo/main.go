package main

import (
	"context"
	"fmt"
	"github-clone/src/model"
	"github-clone/src/repositories"
	"github-clone/src/util"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body string
	statusCode := 201

	user := model.GetUserFromRequest(request)

	repo := model.Repo{}

	repo.FromJSON(request.Body)
	repo.Owner = user

	repository := repositories.GithubRepository{}
	newRepo, err := repository.Create(repo)

	if err != nil {
		httpError := util.HttpErrorFromException(err)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	decodingError, body := newRepo.ToJSON()

	if decodingError.Error == nil {
		statusCode = 201
		fmt.Println("A new repo was created")
	}

	return events.APIGatewayProxyResponse{Body: body, StatusCode: statusCode}, nil
}

func main() {
	lambda.Start(handleRequest)
}

package main

import (
	"context"
	"fmt"
	"github-clone/src/database"
	"github-clone/src/errors"
	"github-clone/src/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body string
	statusCode := 201

	user := model.GetUserFromRequest(request)

	repo := model.Repo{}

	if err := repo.FromJSON(request.Body); err != nil {
		httpError := errors.HttpErrorFromException(err)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	repo.Owner = user

	repository := database.Repository{}
	newRepo, err := repository.Create(repo)

	if err != nil {
		httpError := errors.HttpErrorFromException(err)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	body, decodingError := newRepo.ToJSON()

	if decodingError == nil {
		statusCode = 201
		fmt.Println("A new repo was created")
	}

	return events.APIGatewayProxyResponse{Body: body, StatusCode: statusCode}, nil
}

func main() {
	lambda.Start(handleRequest)
}

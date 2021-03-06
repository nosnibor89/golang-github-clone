package main

import (
	"context"
	"github-clone/src/database/repository"
	"github-clone/src/errors"
	"github-clone/src/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

var (
	repo = model.Repo{}
)

func handleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body string
	statusCode := http.StatusCreated

	user := model.GetUserFromRequest(request)

	if err := repo.FromJSON(request.Body); err != nil {
		httpError := errors.HttpErrorFromException(err)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	repo.Owner = user
	newRepo, err := repository.Create(repo)

	if err != nil {
		httpError := errors.HttpErrorFromException(err)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	body, decodingError := newRepo.ToJSON()

	if decodingError != nil {
		httpError := errors.HttpErrorFromException(decodingError)
		return events.APIGatewayProxyResponse{Body: httpError.Message, StatusCode: httpError.Code}, nil
	}

	return events.APIGatewayProxyResponse{Body: body, StatusCode: statusCode}, nil
}

func main() {
	lambda.Start(handleRequest)
}

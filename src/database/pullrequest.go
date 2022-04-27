package database

import (
	"fmt"
	"github-clone/src/database/entities"
	"github-clone/src/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strconv"
)

type PullRequest struct{}

type PullRequestFindOneInput struct {
	Repo              string
	Owner             string
	PullRequestNumber string
}

func (pr PullRequest) Create(newPullRequest model.PullRequest) (*model.PullRequest, error) {
	pullRequestNumber, err := getNextNumberFromRepo(newPullRequest.Repo.Name, newPullRequest.Repo.Owner.Username)
	if err != nil {
		return nil, err
	}
	newPullRequest.PullRequestNumber = pullRequestNumber

	return pr.createPullRequest(newPullRequest)
}

func (pr PullRequest) createPullRequest(newPullRequest model.PullRequest) (*model.PullRequest, error) {
	prEntity := entities.NewPullRequest(
		newPullRequest.Title,
		newPullRequest.Content,
		newPullRequest.Repo.Name,
		newPullRequest.Repo.Owner.Username,
		newPullRequest.Creator.Username,
		newPullRequest.PullRequestNumber,
	)

	item, err := prEntity.ToItem()

	if err != nil {
		return nil, fmt.Errorf("item resolution errors: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName:           tableName(),
		Item:                item,
		ReturnValues:        aws.String(dynamodb.ReturnValueNone),
		ConditionExpression: generateAttrNotExistsExpression("PK", "SK"),
	}

	if _, err = dynamoDbClient.PutItem(input); err != nil {
		return nil, err
	}

	created := prEntity.ToModel()

	return &created, nil
}

func (pr PullRequest) FindOne(input PullRequestFindOneInput) (*model.PullRequest, error) {
	prNumber, err := strconv.Atoi(input.PullRequestNumber)
	if err != nil {
		return nil, fmt.Errorf("could find pull request %w", err)
	}

	item := entities.PullRequest{
		RepoOwner:         input.Owner,
		RepoName:          input.Repo,
		PullRequestNumber: prNumber,
	}

	getInput := &dynamodb.GetItemInput{
		TableName: tableName(),
		Key:       item.Key(),
	}

	itemOutput, err := dynamoDbClient.GetItem(getInput)

	if err != nil || itemOutput.Item == nil {
		return nil, fmt.Errorf("could not find pull request. error %v, item:: %v\n", err, itemOutput.Item)
	}

	prValue := item.ToModelFromAttrs(itemOutput.Item)

	return &prValue, nil
}

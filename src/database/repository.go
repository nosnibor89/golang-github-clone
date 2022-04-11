package database

import (
	"fmt"
	"github-clone/src/database/entities"
	"github-clone/src/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	dynamoDbClient *dynamodb.DynamoDB
)

type Repository struct{}

func (repo Repository) FindOne(name, owner string) *model.Repo {
	item := entities.GithubRepo{
		Name:  name,
		Owner: owner,
	}

	input := &dynamodb.GetItemInput{
		TableName: tableName(),
		Key:       item.Key(),
	}

	itemOutput, err := dynamoDbClient.GetItem(input)

	if err != nil || itemOutput.Item == nil {
		fmt.Printf("could not find repo. error %v, item:: %v\n", err, itemOutput.Item)
		return nil
	}

	repoValue := item.ToModelFromAttrs(itemOutput.Item)
	return &repoValue
}

func (repo Repository) Create(newRepo model.Repo) (*model.Repo, error) {
	repoEntity := entities.NewGithubRepo(newRepo.Name, newRepo.Owner.Username, newRepo.Description)

	item, err := repoEntity.ToItem()

	if err != nil {
		return nil, fmt.Errorf("item resolution error: %v", err)
	}

	params := &dynamodb.PutItemInput{
		TableName:           tableName(),
		Item:                item,
		ReturnValues:        aws.String(dynamodb.ReturnValueNone),
		ConditionExpression: generateAttrNotExistsExpression("PK"),
	}

	if _, err = dynamoDbClient.PutItem(params); err != nil {
		return nil, err
	}

	created := repoEntity.ToModel()

	return &created, nil
}

func (repo Repository) Delete(name, owner string) error {
	repoEntity := entities.NewGithubRepo(name, owner, "")

	params := &dynamodb.DeleteItemInput{
		TableName:    tableName(),
		Key:          repoEntity.Key(),
		ReturnValues: aws.String(dynamodb.ReturnValueNone),
	}

	if _, err := dynamoDbClient.DeleteItem(params); err != nil {
		return err
	}

	return nil
}

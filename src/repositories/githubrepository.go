package repositories

import (
	"fmt"
	"github-clone/src/model"
	"github-clone/src/repositories/entities"
	"github-clone/src/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
	"time"
)

var (
	dynamoDbClient *dynamodb.DynamoDB
)

const (
	tableNameEnvVar = "GITHUB_TABLE_NAME"
)

type GithubRepository struct{}

func (repo GithubRepository) tableName() *string {
	return aws.String(os.Getenv(tableNameEnvVar))
}

func (repo GithubRepository) FindOne(name, owner string) model.Repo {
	entity := entities.GithubRepo{
		Name:  name,
		Owner: owner,
	}

	input := &dynamodb.GetItemInput{
		TableName: repo.tableName(),
		Key:       entity.Key(),
	}

	itemOutput, err := dynamoDbClient.GetItem(input)

	if err != nil {
		fmt.Println("Error occurred find one repo", err)
		return model.Repo{}
	}

	if itemOutput.Item == nil {
		return model.Repo{}
	}

	return entity.ToModelFromAttrs(itemOutput.Item)
}

func (repo GithubRepository) Create(newRepo model.Repo) (model.Repo, error) {
	repoEntity := entities.GithubRepo{
		Name:        newRepo.Name,
		Owner:       newRepo.Owner.Username,
		Description: newRepo.Description,
		UpdatedAt:   time.Now(),
		CreatedAt:   time.Now(),
	}

	item, err := repoEntity.ToItem()

	if err != nil {
		return model.Repo{}, nil
	}

	params := &dynamodb.PutItemInput{
		TableName:           repo.tableName(),
		Item:                item,
		ReturnValues:        aws.String(dynamodb.ReturnValueNone),
		ConditionExpression: util.GenerateAttrNotExistsExpression("PK"),
	}

	if _, err = dynamoDbClient.PutItem(params); err != nil {
		return model.Repo{}, err
	}

	return repoEntity.ToModel(), nil
}

func (repo GithubRepository) Delete(name, owner string) error {
	entity := entities.GithubRepo{
		Name:  name,
		Owner: owner,
	}

	params := &dynamodb.DeleteItemInput{
		TableName:    repo.tableName(),
		Key:          entity.Key(),
		ReturnValues: aws.String(dynamodb.ReturnValueNone),
	}

	if _, err := dynamoDbClient.DeleteItem(params); err != nil {
		return err
	}

	return nil
}

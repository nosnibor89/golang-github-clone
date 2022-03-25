package repositories

import (
	"fmt"
	"github-clone/src/model"
	"github-clone/src/repositories/entities"
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

func (repo GithubRepository) Create(newRepo model.Repo) (model.Repo, error) {
	repoEntity := entities.GithubRepo{
		Name:      newRepo.Name,
		Owner:     newRepo.Owner.Username,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	item, err := repoEntity.ToItem()

	if err != nil {
		return model.Repo{}, nil
	}

	fmt.Println("tableNameEnvVar")
	fmt.Println(os.Getenv(tableNameEnvVar))

	params := &dynamodb.PutItemInput{
		TableName:    aws.String(os.Getenv(tableNameEnvVar)),
		Item:         item,
		ReturnValues: aws.String(dynamodb.ReturnValueNone),
	}

	putItemOutput, err := dynamoDbClient.PutItem(params)

	if err != nil {
		return model.Repo{}, err
	}

	return repoEntity.ToModel(putItemOutput.Attributes), nil
}

func (repo GithubRepository) Delete(name, owner string) error {
	entity := entities.GithubRepo{
		Name:  name,
		Owner: owner,
	}

	params := &dynamodb.DeleteItemInput{
		TableName:    aws.String(os.Getenv(tableNameEnvVar)),
		Key:          entity.Key(),
		ReturnValues: aws.String(dynamodb.ReturnValueNone),
	}

	if _, err := dynamoDbClient.DeleteItem(params); err != nil {
		return err
	}

	return nil
}

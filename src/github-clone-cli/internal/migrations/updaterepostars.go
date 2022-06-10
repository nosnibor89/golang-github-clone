package migrations

import (
	"fmt"
	"github-clone/src/database"
	"github-clone/src/database/internal/entities"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type UpdateRepoStarsMigration struct {
}

func (m UpdateRepoStarsMigration) Name() string {
	return "update-repo-stars"
}

func (m UpdateRepoStarsMigration) Run() error {
	return database.ScanFilterByType("REPO", updateRepo)
}

func updateRepo(item map[string]*dynamodb.AttributeValue) error {
	fmt.Printf("Updating REPO PK %s, SK %s \n", aws.StringValue(item["PK"].S), aws.StringValue(item["SK"].S))
	repoEntity := entities.GithubRepo{
		Name:  aws.StringValue(item["Name"].S),
		Owner: aws.StringValue(item["Owner"].S),
	}

	input := &dynamodb.UpdateItemInput{
		TableName:        database.TableName(),
		Key:              repoEntity.Key(),
		UpdateExpression: aws.String("ADD StarCount :count REMOVE StartCount"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":count": {
				N: aws.String("0"),
			},
		},
	}

	if _, err := database.DBClient().UpdateItem(input); err != nil {
		return err
	}

	return nil
}

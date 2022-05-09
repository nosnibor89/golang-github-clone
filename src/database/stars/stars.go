package stars

import (
	"fmt"
	"github-clone/src/database"
	"github-clone/src/database/entities"
	"github-clone/src/model"
	"github-clone/src/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
)

func StarRepo(repo, owner, username string) error {
	repoItem := entities.GithubRepo{
		Name:  repo,
		Owner: owner,
	}

	star := entities.NewStar(repo, owner, username)

	starItem, err := star.ToItem()

	if err != nil {
		return fmt.Errorf("item resolution error: %v", err)
	}

	//TODO: Add error handling
	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					TableName:           database.TableName(),
					Item:                starItem,
					ConditionExpression: database.GenerateAttrNotExistsExpression("PK", "SK"),
				},
			},
			{
				Update: &dynamodb.Update{
					TableName:        database.TableName(),
					Key:              repoItem.Key(),
					UpdateExpression: aws.String("SET #startCount = #startCount + :incr"),
					ExpressionAttributeNames: map[string]*string{
						"#startCount": aws.String("StarCount"),
					},
					ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
						":incr": {
							N: aws.String("1"),
						},
					},
				},
			},
		},
	}

	if _, err := database.DBClient().TransactWriteItems(input); err != nil {
		return fmt.Errorf("could not star repo: %w", err)
	}

	return nil
}

func UnStarRepo(repo, owner, username string) error {
	repoItem := entities.GithubRepo{
		Name:  repo,
		Owner: owner,
	}

	star := entities.Star{
		RepoName:  repo,
		RepoOwner: owner,
		Username:  username,
	}

	//TODO: Add error handling
	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Delete: &dynamodb.Delete{
					TableName: database.TableName(),
					Key:       star.Key(),
				},
			},
			{
				Update: &dynamodb.Update{
					TableName:        database.TableName(),
					Key:              repoItem.Key(),
					UpdateExpression: aws.String("SET #startCount = #startCount - :incr"),
					ExpressionAttributeNames: map[string]*string{
						"#startCount": aws.String("StarCount"),
					},
					ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
						":incr": {
							N: aws.String("1"),
						},
					},
				},
			},
		},
	}

	if _, err := database.DBClient().TransactWriteItems(input); err != nil {
		return fmt.Errorf("could not un-star repo: %w", err)
	}

	return nil
}

func FindStargazers(repo, owner string) (*model.Repo, []string, error) {
	var stargazers []string

	entity := entities.Star{
		RepoOwner: owner,
		RepoName:  repo,
	}

	pk := entity.PartitionKey()

	input := &dynamodb.QueryInput{
		TableName:              database.TableName(),
		KeyConditionExpression: aws.String("#pk = :pk AND #sk >= :sk"),
		ExpressionAttributeNames: map[string]*string{
			"#pk": aws.String("PK"),
			"#sk": aws.String("SK"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(pk),
			},
			":sk": {
				S: aws.String(pk),
			},
		},
		ReturnConsumedCapacity: aws.String(dynamodb.ReturnConsumedCapacityTotal),
	}

	queryOutput, err := database.DBClient().Query(input)

	if err != nil {
		return nil, stargazers, fmt.Errorf("error fetching stargazers: %w", err)
	}

	repoItem, starItems := queryOutput.Items[0], queryOutput.Items[1:]

	if *queryOutput.ScannedCount > *queryOutput.Count {
		util.LogYellow("WARNING: ScannedCount is greater than Count")
		log.Printf("[Trace]ScannedCount: %d", *queryOutput.ScannedCount)
		log.Printf("[Trace]Count: %d", *queryOutput.Count)
	}

	repoEntity := entities.GithubRepo{}
	repoValue := repoEntity.ToModelFromAttrs(repoItem)
	return &repoValue, toStargazers(starItems), nil
}

func toStargazers(starItems []entities.Attrs) []string {
	var stargazers []string

	for _, item := range starItems {
		stargazers = append(stargazers, aws.StringValue(item["Username"].S))
	}

	return stargazers
}

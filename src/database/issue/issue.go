package issue

import (
	"fmt"
	"github-clone/src/database"
	entities2 "github-clone/src/database/internal/entities"
	"github-clone/src/database/repository"
	"github-clone/src/model"
	"github-clone/src/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"strconv"
	"strings"
)

type FindIssueInput struct {
	Repo, Owner, IssueNumber string
}

func Create(newIssue model.Issue) (*model.Issue, error) {
	issueNumber, err := repository.GetNextNumberFromRepo(newIssue.Repo.Name, newIssue.Repo.Owner.Username)
	if err != nil {
		return nil, err
	}

	newIssue.IssueNumber = issueNumber
	return createIssue(newIssue)
}

func Find(repo, owner, status string) (*model.Repo, []model.Issue, error) {
	var issues []model.Issue
	shouldLookOpenIssues := status == "" || strings.ToUpper(strings.TrimSpace(status)) == entities2.IssueOpenStatus

	entity := entities2.Issue{
		RepoOwner: owner,
		RepoName:  repo,
	}

	input := &dynamodb.QueryInput{
		TableName:              database.TableName(),
		KeyConditionExpression: aws.String("#pk = :pk AND #sk <= :sk"),
		ExpressionAttributeNames: map[string]*string{
			"#pk":   aws.String("PK"),
			"#sk":   aws.String("SK"),
			"#open": aws.String("Open"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(entity.PartitionKey()),
			},
			":sk": {
				S: aws.String(entity.PartitionKey()),
			},
			":open": {
				BOOL: aws.Bool(shouldLookOpenIssues),
			},
		},
		ScanIndexForward:       aws.Bool(false),
		FilterExpression:       aws.String("(attribute_not_exists(#open)) OR (#open = :open)"),
		ReturnConsumedCapacity: aws.String(dynamodb.ReturnConsumedCapacityTotal),
	}

	queryOutput, err := database.DBClient().Query(input)

	if err != nil {
		return nil, issues, fmt.Errorf("error fetching issues: %w", err)
	}

	repoItem, issueItems := queryOutput.Items[0], queryOutput.Items[1:]

	if *queryOutput.ScannedCount > *queryOutput.Count {
		util.LogYellow("WARNING: ScannedCount is greater than Count")
		log.Printf("[Trace]ScannedCount: %d", *queryOutput.ScannedCount)
		log.Printf("[Trace]Count: %d", *queryOutput.Count)
	}

	issues = entities2.ToList[model.Issue](issueItems, entity)

	repoEntity := entities2.GithubRepo{}

	repoValue := repoEntity.ToModelFromAttrs(repoItem)
	return &repoValue, issues, nil
}

func FindOne(input FindIssueInput) (*model.Issue, error) {
	issueNumber, err := strconv.Atoi(input.IssueNumber)
	if err != nil {
		return nil, fmt.Errorf("could find issue %w", err)
	}

	item := entities2.Issue{
		IssueNumber: issueNumber,
		RepoName:    input.Repo,
		RepoOwner:   input.Owner,
	}

	dynamoInput := &dynamodb.GetItemInput{
		TableName: database.TableName(),
		Key:       item.Key(),
	}

	itemOutput, err := database.DBClient().GetItem(dynamoInput)

	if err != nil || itemOutput.Item == nil {
		return nil, fmt.Errorf("could not find issue. error %v, item:: %v\n", err, itemOutput.Item)
	}

	issueValue := item.ToModelFromAttrs(itemOutput.Item)

	return &issueValue, nil
}

func createIssue(newIssue model.Issue) (*model.Issue, error) {
	issueEntity := entities2.NewIssue(
		newIssue.Title,
		newIssue.Content,
		newIssue.Repo.Name,
		newIssue.Repo.Owner.Username,
		newIssue.Creator.Username,
		newIssue.IssueNumber,
	)

	item, err := issueEntity.ToItem()

	if err != nil {
		return nil, fmt.Errorf("item resolution error: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName:           database.TableName(),
		Item:                item,
		ReturnValues:        aws.String(dynamodb.ReturnValueNone),
		ConditionExpression: database.GenerateAttrNotExistsExpression("PK", "SK"),
	}

	if _, err = database.DBClient().PutItem(input); err != nil {
		return nil, err
	}

	created := issueEntity.ToModel()

	return &created, err
}

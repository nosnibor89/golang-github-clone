package database

import (
	"fmt"
	"github-clone/src/database/entities"
	"github-clone/src/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"strconv"
	"strings"
)

type Issue struct {
}

type IssueFindOneInput struct {
	Repo, Owner, IssueNumber string
}

const noIssueNumber = "could not assign issue number"

func (i Issue) Create(newIssue model.Issue) (*model.Issue, error) {
	issueNumber, err := i.getIssueNumberFromRepo(newIssue.Repo.Name, newIssue.Repo.Owner.Username)
	if err != nil {
		return nil, err
	}

	newIssue.IssueNumber = issueNumber
	return i.createIssue(newIssue)
}

func (i Issue) GetIssues(repo, owner, status string) (*model.Repo, []model.Issue, error) {
	var issues []model.Issue
	shouldLookOpenIssues := status == "" || strings.ToUpper(strings.TrimSpace(status)) == entities.IssueOpenStatus

	entity := entities.Issue{
		RepoOwner: owner,
		RepoName:  repo,
	}

	input := &dynamodb.QueryInput{
		TableName:              tableName(),
		KeyConditionExpression: aws.String("#pk = :pk"),
		ExpressionAttributeNames: map[string]*string{
			"#pk":   aws.String("PK"),
			"#open": aws.String("Open"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
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

	queryOutput, err := dynamoDbClient.Query(input)

	if err != nil {
		return nil, issues, fmt.Errorf("error fetching issues: %w", err)
	}

	repoItem, issueItems := queryOutput.Items[0], queryOutput.Items[1:]

	log.Printf("[Trace]ScannedCount: %d", *queryOutput.ScannedCount)
	log.Printf("[Trace]Count: %d", *queryOutput.Count)
	issues = entity.ToIssueList(issueItems)

	repoEntity := entities.GithubRepo{}

	repoValue := repoEntity.ToModelFromAttrs(repoItem)
	return &repoValue, issues, nil
}

func (i Issue) FindOne(input IssueFindOneInput) (*model.Issue, error) {
	issueNumber, err := strconv.Atoi(input.IssueNumber)
	if err != nil {
		return nil, fmt.Errorf("could find issue %w", err)
	}

	item := entities.Issue{
		IssueNumber: issueNumber,
		RepoName:    input.Repo,
		RepoOwner:   input.Owner,
	}

	dynamoInput := &dynamodb.GetItemInput{
		TableName: tableName(),
		Key:       item.Key(),
	}

	itemOutput, err := dynamoDbClient.GetItem(dynamoInput)

	if err != nil || itemOutput.Item == nil {
		return nil, fmt.Errorf("could not find issue. error %v, item:: %v\n", err, itemOutput.Item)
	}

	issueValue := item.ToModelFromAttrs(itemOutput.Item)

	return &issueValue, nil
}

func (i Issue) createIssue(newIssue model.Issue) (*model.Issue, error) {
	issueEntity := entities.NewIssue(
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
		TableName:           tableName(),
		Item:                item,
		ReturnValues:        aws.String(dynamodb.ReturnValueNone),
		ConditionExpression: generateAttrNotExistsExpression("PK", "SK"),
	}

	if _, err = dynamoDbClient.PutItem(input); err != nil {
		return nil, err
	}

	created := issueEntity.ToModel()

	return &created, err
}

func (i Issue) updateRepo(repo, owner string) (map[string]*dynamodb.AttributeValue, error) {
	repoEntity := entities.GithubRepo{
		Name:  repo,
		Owner: owner,
	}

	expressionAttrs := map[string]*string{
		"#issuePRNumber": aws.String("IssuePRNumber"),
	}

	expressionAttrsValues := map[string]*dynamodb.AttributeValue{
		":incr": {
			N: aws.String("1"),
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 tableName(),
		Key:                       repoEntity.Key(),
		UpdateExpression:          aws.String("SET #issuePRNumber = #issuePRNumber + :incr"),
		ExpressionAttributeNames:  expressionAttrs,
		ExpressionAttributeValues: expressionAttrsValues,
		ReturnValues:              aws.String(dynamodb.ReturnValueAllNew),
	}

	updated, err := dynamoDbClient.UpdateItem(input)

	if err != nil {
		return nil, fmt.Errorf("%s(failed updating repository)", noIssueNumber)
	}

	return updated.Attributes, nil
}

func (i Issue) getIssueNumberFromRepo(repo, owner string) (int, error) {
	var issueNumber int
	updatedAttrs, err := i.updateRepo(repo, owner)
	if err != nil {
		return issueNumber, err
	}

	issueNumberAttr := updatedAttrs["IssuePRNumber"]
	if issueNumberAttr == nil {
		return issueNumber, fmt.Errorf("%s(issue number is not set in repository)", noIssueNumber)
	}

	issueNumber, err = strconv.Atoi(aws.StringValue(issueNumberAttr.N))

	if err != nil {
		return 0, fmt.Errorf("%v(issue number could not be parsed)::::%v", noIssueNumber, err)
	}

	return issueNumber, nil
}

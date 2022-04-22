package database

import (
	"fmt"
	"github-clone/src/database/entities"
	"github-clone/src/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strconv"
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

func getNextNumberFromRepo(repo, owner string) (int, error) {
	var issueNumber int
	updatedAttrs, err := updateRepo(repo, owner)
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

func updateRepo(repo, owner string) (map[string]*dynamodb.AttributeValue, error) {
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
		return nil, fmt.Errorf("%s(failed updating repository)\n %w", noIssueNumber, err)
	}

	return updated.Attributes, nil
}

package repository

import (
	"fmt"
	"github-clone/src/database"
	"github-clone/src/database/entities"
	"github-clone/src/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strconv"
)

const noIssuePRNumberMessage = "could not assign issue number"

func FindRepository(name, owner string) *model.Repo {
	item := entities.GithubRepo{
		Name:  name,
		Owner: owner,
	}

	input := &dynamodb.GetItemInput{
		TableName: database.TableName(),
		Key:       item.Key(),
	}

	itemOutput, err := database.DBClient().GetItem(input)

	if err != nil || itemOutput.Item == nil {
		fmt.Printf("could not find repo. error %v, item:: %v\n", err, itemOutput.Item)
		return nil
	}

	repoValue := item.ToModelFromAttrs(itemOutput.Item)
	return &repoValue
}

func CreateRepository(newRepo model.Repo) (*model.Repo, error) {
	repoEntity := entities.NewGithubRepo(newRepo.Name, newRepo.Owner.Username, newRepo.Description)

	item, err := repoEntity.ToItem()

	if err != nil {
		return nil, fmt.Errorf("item resolution error: %v", err)
	}

	params := &dynamodb.PutItemInput{
		TableName:           database.TableName(),
		Item:                item,
		ReturnValues:        aws.String(dynamodb.ReturnValueNone),
		ConditionExpression: database.GenerateAttrNotExistsExpression("PK"),
	}

	if _, err = database.DBClient().PutItem(params); err != nil {
		return nil, err
	}

	created := repoEntity.ToModel()

	return &created, nil
}

func DeleteRepository(name, owner string) error {
	repoEntity := entities.NewGithubRepo(name, owner, "")

	params := &dynamodb.DeleteItemInput{
		TableName:    database.TableName(),
		Key:          repoEntity.Key(),
		ReturnValues: aws.String(dynamodb.ReturnValueNone),
	}

	if _, err := database.DBClient().DeleteItem(params); err != nil {
		return err
	}

	return nil
}

func GetNextNumberFromRepo(repo, owner string) (int, error) {
	var number int
	updatedAttrs, err := updateRepo(repo, owner)
	if err != nil {
		return number, err
	}

	numberAttr := updatedAttrs["IssuePRNumber"]
	if numberAttr == nil {
		return number, fmt.Errorf("%s(IssuePRNumber number is not set in repository)", noIssuePRNumberMessage)
	}

	number, err = strconv.Atoi(aws.StringValue(numberAttr.N))

	if err != nil {
		return 0, fmt.Errorf("%v(issue number could not be parsed)::::%v", noIssuePRNumberMessage, err)
	}

	return number, nil
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
		TableName:                 database.TableName(),
		Key:                       repoEntity.Key(),
		UpdateExpression:          aws.String("SET #issuePRNumber = #issuePRNumber + :incr"),
		ExpressionAttributeNames:  expressionAttrs,
		ExpressionAttributeValues: expressionAttrsValues,
		ReturnValues:              aws.String(dynamodb.ReturnValueAllNew),
	}

	updated, err := database.DBClient().UpdateItem(input)

	if err != nil {
		return nil, fmt.Errorf("%s(failed updating repository)\n %w", noIssuePRNumberMessage, err)
	}

	return updated.Attributes, nil
}

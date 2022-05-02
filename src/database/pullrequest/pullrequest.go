package pullrequest

import (
	"fmt"
	"github-clone/src/database"
	"github-clone/src/database/entities"
	"github-clone/src/database/repository"
	"github-clone/src/model"
	"github-clone/src/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"strconv"
	"strings"
)

type FindPullRequestInput struct {
	Repo              string
	Owner             string
	PullRequestNumber string
}

func CreatePullRequest(newPullRequest model.PullRequest) (*model.PullRequest, error) {
	pullRequestNumber, err := repository.GetNextNumberFromRepo(newPullRequest.Repo.Name, newPullRequest.Repo.Owner.Username)
	if err != nil {
		return nil, err
	}
	newPullRequest.PullRequestNumber = pullRequestNumber

	return createPullRequest(newPullRequest)
}

func createPullRequest(newPullRequest model.PullRequest) (*model.PullRequest, error) {
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
		TableName:           database.TableName(),
		Item:                item,
		ReturnValues:        aws.String(dynamodb.ReturnValueNone),
		ConditionExpression: database.GenerateAttrNotExistsExpression("PK", "SK"),
	}

	if _, err = database.DBClient().PutItem(input); err != nil {
		return nil, err
	}

	created := prEntity.ToModel()

	return &created, nil
}

func FindPullRequest(input FindPullRequestInput) (*model.PullRequest, error) {
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
		TableName: database.TableName(),
		Key:       item.Key(),
	}

	itemOutput, err := database.DBClient().GetItem(getInput)

	if err != nil || itemOutput.Item == nil {
		return nil, fmt.Errorf("could not find pull request. error %v, item:: %v\n", err, itemOutput.Item)
	}

	prValue := item.ToModelFromAttrs(itemOutput.Item)

	return &prValue, nil
}

func GetPullRequests(repo, owner, status string) (*model.Repo, []model.PullRequest, error) {
	var prs []model.PullRequest
	shouldLookOpen := status == "" || strings.ToUpper(strings.TrimSpace(status)) == entities.PullRequestOpenStatus

	entity := entities.PullRequest{
		RepoOwner: owner,
		RepoName:  repo,
	}

	gsiPK, gsiSK := entity.GSI1()

	input := &dynamodb.QueryInput{
		TableName:              database.TableName(),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("#pk = :pk AND #sk >= :sk"),
		ExpressionAttributeNames: map[string]*string{
			"#pk":   aws.String("GSI1PK"),
			"#sk":   aws.String("GSI1SK"),
			"#open": aws.String("Open"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(gsiPK),
			},
			":sk": {
				S: aws.String(gsiSK),
			},
			":open": {
				BOOL: aws.Bool(shouldLookOpen),
			},
		},
		ScanIndexForward:       aws.Bool(false),
		FilterExpression:       aws.String("(attribute_not_exists(#open)) OR (#open = :open)"),
		ReturnConsumedCapacity: aws.String(dynamodb.ReturnConsumedCapacityTotal),
	}

	queryOutput, err := database.DBClient().Query(input)

	if err != nil {
		return nil, prs, fmt.Errorf("error fetching pull requests: %w", err)
	}

	repoItem, prItems := queryOutput.Items[0], queryOutput.Items[1:]

	if *queryOutput.ScannedCount > *queryOutput.Count {
		util.LogYellow("WARNING: ScannedCount is greater than Count")
		log.Printf("[Trace]ScannedCount: %d", *queryOutput.ScannedCount)
		log.Printf("[Trace]Count: %d", *queryOutput.Count)
	}

	prs = entities.ToList[model.PullRequest](prItems, entity)
	repoEntity := entities.GithubRepo{}

	repoValue := repoEntity.ToModelFromAttrs(repoItem)
	return &repoValue, prs, nil
}

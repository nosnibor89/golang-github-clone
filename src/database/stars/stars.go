package stars

import (
	"fmt"
	"github-clone/src/database"
	"github-clone/src/database/entities"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

	//nameBuilder := expression.Name("StartCount")
	//operandBuild := expression.Plus(nameBuilder, expression.Value("1"))
	//updateExpression := expression.Set(nameBuilder, operandBuild)
	//expr, err := expression.NewBuilder().WithUpdate(updateExpression).Build()
	//if err != nil {
	//	fmt.Println(err)
	//}

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
					TableName: database.TableName(),
					Key:       repoItem.Key(),
					//UpdateExpression:          expr.Update(),
					//ExpressionAttributeNames:  expr.Names(),
					//ExpressionAttributeValues: expr.Values(),
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

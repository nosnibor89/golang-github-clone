package entities

import (
	"fmt"
	"github-clone/src/model"
	"github.com/aws/aws-sdk-go/aws"
	"strconv"
	"strings"
	"time"
)

type Issue struct {
	Entity
	Open                                         bool
	RepoName, RepoOwner, Content, Title, Creator string
	IssueNumber                                  int
}

const (
	IssueOpenStatus = "OPEN"
	//IssueClosedStatus = "CLOSED"
)

func NewIssue(title, content, repoName, repoOwner, creator string, issueNumber int) Issue {
	return Issue{
		Title:       title,
		Content:     content,
		RepoName:    repoName,
		RepoOwner:   repoOwner,
		IssueNumber: issueNumber,
		Creator:     creator,
		Open:        true,
		Entity: Entity{
			time.Now(),
			time.Now(),
		},
	}
}

func (i Issue) PartitionKey() string {
	return fmt.Sprintf("REPO#%s#%s", strings.ToLower(i.RepoOwner), strings.ToLower(i.RepoName))
}

func (i Issue) Key() Attrs {
	ad := i.initialAttributeDefinition()

	return ad.getPrimaryKey()
}

// ToItem Exports entity to Item type
func (i Issue) ToItem() (Attrs, error) {
	//TODO: ADD VALIDATION

	ad := i.initialAttributeDefinition()
	ad.withStringAttribute("Title", i.Title).
		withStringAttribute("Content", i.Content).
		withStringAttribute("Creator", i.Creator).
		withStringAttribute("RepoOwner", i.RepoOwner).
		withStringAttribute("RepoName", i.RepoName).
		withBoolAttribute("Open", i.Open).
		withStringAttribute("CreatedAt", parseTimeItem(i.CreatedAt)).
		withStringAttribute("UpdatedAt", parseTimeItem(i.UpdatedAt)).
		withIntAttribute("IssueNumber", strconv.Itoa(i.IssueNumber))

	itemAttrs := ad.allAttributes()

	return itemAttrs, nil
}

func (i Issue) ToModelFromAttrs(attrs Attrs) model.Issue {
	var repoModel model.Issue
	if attrs["Title"] != nil && attrs["Creator"] != nil {
		creator := model.User{Username: aws.StringValue(attrs["Creator"].S), Name: aws.StringValue(attrs["Creator"].S)}

		issueNumber, _ := strconv.Atoi(aws.StringValue(attrs["IssueNumber"].N))

		repoModel = model.Issue{
			Model: model.Model{
				UpdatedAt: parseTimeAttr(aws.StringValue(attrs["UpdatedAt"].S)),
				CreatedAt: parseTimeAttr(aws.StringValue(attrs["CreatedAt"].S)),
			},
			Title:       aws.StringValue(attrs["Title"].S),
			Content:     aws.StringValue(attrs["Content"].S),
			IssueNumber: issueNumber,
			Open:        aws.BoolValue(attrs["Open"].BOOL),
			Creator:     creator,
		}

	}
	return repoModel
}

func (i Issue) ToModel() model.Issue {
	repoModel := model.Issue{
		Model: model.Model{
			UpdatedAt: i.UpdatedAt,
			CreatedAt: i.CreatedAt,
		},
		Title:       i.Title,
		Content:     i.Content,
		IssueNumber: i.IssueNumber,
		Open:        i.Open,
		Creator:     model.User{Username: i.Creator, Name: i.Creator},
	}
	return repoModel
}

func pad(num int) string {
	return fmt.Sprintf("%07d", num)
}

func (i Issue) initialAttributeDefinition() attrDefinition {
	sk := fmt.Sprintf("ISSUE#%s", pad(i.IssueNumber))
	return attrDefinition{
		pk:        i.PartitionKey(),
		sk:        sk,
		typeLabel: "ISSUE",
	}
}

func (i Issue) ToIssueList(items []Attrs) []model.Issue {
	issues := []model.Issue{}

	for _, item := range items {
		issues = append(issues, i.ToModelFromAttrs(item))
	}

	return issues
}

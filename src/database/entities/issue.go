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
	RepoOwner, Content, Title, Creator string
	IssueNumber                        int
}

func NewIssue(title, content, repoOwner, creator string, issueNumber int) Issue {
	return Issue{
		Title:       title,
		Content:     content,
		RepoOwner:   repoOwner,
		IssueNumber: issueNumber,
		Creator:     creator,
		Entity: Entity{
			time.Now(),
			time.Now(),
		},
	}
}

func (issue Issue) Key() Attrs {
	ad := issue.initialAttributeDefinition(issue.RepoOwner, issue.IssueNumber)

	return ad.getPrimaryKey()
}

// ToItem Exports entity to Item type
func (issue Issue) ToItem() (Attrs, error) {
	//TODO: ADD VALIDATION

	ad := issue.initialAttributeDefinition(issue.RepoOwner, issue.IssueNumber)
	ad.withStringAttribute("Title", issue.Title).
		withStringAttribute("Content", issue.Content).
		withStringAttribute("Creator", issue.Creator).
		withStringAttribute("RepoOwner", issue.RepoOwner).
		withIntAttribute("IssueNumber", strconv.Itoa(issue.IssueNumber))

	itemAttrs := ad.allAttributes()

	return itemAttrs, nil
}

func (issue Issue) ToModelFromAttrs(attrs Attrs) model.Issue {
	var repoModel model.Issue
	if attrs["Title"] != nil && attrs["Creator"] != nil {
		creator := model.User{Username: aws.StringValue(attrs["Creator"].S), Name: aws.StringValue(attrs["Creator"].S)}

		repoModel = model.Issue{
			Model: model.Model{
				Identifier: aws.StringValue(attrs["Title"].S),
			},
			Content:   aws.StringValue(attrs["Content"].S),
			Creator:   creator,
			UpdatedAt: parseTimeAttr(aws.StringValue(attrs["UpdatedAt"].S)),
			CreatedAt: parseTimeAttr(aws.StringValue(attrs["CreatedAt"].S)),
		}

	}
	return repoModel
}

func (issue Issue) ToModel() model.Issue {
	repoModel := model.Issue{
		Model: model.Model{
			Identifier: issue.Title,
		},
		Content:   issue.Content,
		Title:     issue.Title,
		Creator:   model.User{Username: issue.Creator, Name: issue.Creator},
		UpdatedAt: issue.UpdatedAt,
		CreatedAt: issue.CreatedAt,
	}
	return repoModel
}

func pad(num int) string {
	return fmt.Sprintf("%07d", num)
}

func (issue Issue) initialAttributeDefinition(owner string, issueNumber int) attrDefinition {
	pk := fmt.Sprintf("REPO#%s#%s", strings.ToLower(owner), strings.ToLower(owner))
	sk := fmt.Sprintf("ISSUE#%s", pad(issueNumber))
	return attrDefinition{
		pk:        pk,
		sk:        sk,
		typeLabel: "ISSUE",
	}
}

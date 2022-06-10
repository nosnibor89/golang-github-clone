package entities

import (
	"fmt"
	"github-clone/src/model"
	"github.com/aws/aws-sdk-go/aws"
	"strconv"
	"strings"
	"time"
)

//TODO: Maybe avoid duplication ??
type PullRequest struct {
	Entity
	Open                                         bool
	RepoName, RepoOwner, Content, Title, Creator string
	PullRequestNumber                            int
}

const (
	PullRequestOpenStatus = "OPEN"
)

func NewPullRequest(title, content, repoName, repoOwner, creator string, pullRequestNumber int) PullRequest {
	return PullRequest{
		Title:             title,
		Content:           content,
		RepoName:          repoName,
		RepoOwner:         repoOwner,
		PullRequestNumber: pullRequestNumber,
		Creator:           creator,
		Open:              true,
		Entity: Entity{
			time.Now(),
			time.Now(),
		},
	}
}

// ToItem Exports entity to Item type
func (pr PullRequest) ToItem() (Attrs, error) {
	//TODO: ADD VALIDATION
	gs1pk, gs1sk := pr.GSI1()

	ad := pr.initialAttributeDefinition()
	ad.withStringAttribute("Title", pr.Title).
		withStringAttribute("Content", pr.Content).
		withStringAttribute("Creator", pr.Creator).
		withStringAttribute("RepoOwner", pr.RepoOwner).
		withStringAttribute("RepoName", pr.RepoName).
		withBoolAttribute("Open", pr.Open).
		withStringAttribute("CreatedAt", parseTimeItem(pr.CreatedAt)).
		withStringAttribute("UpdatedAt", parseTimeItem(pr.UpdatedAt)).
		withIntAttribute("PullRequestNumber", strconv.Itoa(pr.PullRequestNumber)).
		withSecondaryIndexKey(1, gs1pk, gs1sk)

	itemAttrs := ad.allAttributes()

	return itemAttrs, nil
}

func (pr PullRequest) GSI1() (string, string) {
	gs1pk := fmt.Sprintf("REPO#%s#%s", strings.ToLower(pr.RepoOwner), strings.ToLower(pr.RepoName))
	gs1sk := fmt.Sprintf("PR#%s", pad(pr.PullRequestNumber))

	return gs1pk, gs1sk
}

func (pr PullRequest) PartitionKey() string {
	return fmt.Sprintf("PR#%s#%s#%s", strings.ToLower(pr.RepoOwner), strings.ToLower(pr.RepoName), strconv.Itoa(pr.PullRequestNumber))
}

func (pr PullRequest) Key() Attrs {
	ad := pr.initialAttributeDefinition()

	return ad.getPrimaryKey()
}

func (pr PullRequest) initialAttributeDefinition() attrDefinition {
	sk := fmt.Sprintf("PR#%s#%s#%s", strings.ToLower(pr.RepoOwner), strings.ToLower(pr.RepoName), pad(pr.PullRequestNumber))
	return attrDefinition{
		pk:        pr.PartitionKey(),
		sk:        sk,
		typeLabel: "PULL_REQUEST",
	}
}

func (pr PullRequest) ToModelFromAttrs(attrs Attrs) model.PullRequest {
	var pullRequest model.PullRequest
	if attrs["Title"] != nil && attrs["Creator"] != nil {
		creator := model.User{Username: aws.StringValue(attrs["Creator"].S), Name: aws.StringValue(attrs["Creator"].S)}

		prNumber, _ := strconv.Atoi(aws.StringValue(attrs["PullRequestNumber"].N))

		pullRequest = model.PullRequest{
			PullRequestNumber: prNumber,
			Model: model.Model{
				UpdatedAt: parseTimeAttr(aws.StringValue(attrs["UpdatedAt"].S)),
				CreatedAt: parseTimeAttr(aws.StringValue(attrs["CreatedAt"].S)),
			},
			RepoElement: model.RepoElement{
				Title:   aws.StringValue(attrs["Title"].S),
				Content: aws.StringValue(attrs["Content"].S),
				Open:    aws.BoolValue(attrs["Open"].BOOL),
				Creator: creator,
			},
		}

	}
	return pullRequest
}

func (pr PullRequest) ToModel() model.PullRequest {
	prModel := model.PullRequest{
		Model: model.Model{
			UpdatedAt: pr.UpdatedAt,
			CreatedAt: pr.CreatedAt,
		},
		PullRequestNumber: pr.PullRequestNumber,
		RepoElement: model.RepoElement{
			Title:   pr.Title,
			Content: pr.Content,
			Open:    pr.Open,
			Creator: model.User{Username: pr.Creator, Name: pr.Creator},
		},
	}
	return prModel
}

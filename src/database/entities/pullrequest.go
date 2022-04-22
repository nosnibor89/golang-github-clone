package entities

import (
	"fmt"
	"github-clone/src/model"
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

	ad := pr.initialAttributeDefinition()
	ad.withStringAttribute("Title", pr.Title).
		withStringAttribute("Content", pr.Content).
		withStringAttribute("Creator", pr.Creator).
		withStringAttribute("RepoOwner", pr.RepoOwner).
		withStringAttribute("RepoName", pr.RepoName).
		withBoolAttribute("Open", pr.Open).
		withStringAttribute("CreatedAt", parseTimeItem(pr.CreatedAt)).
		withStringAttribute("UpdatedAt", parseTimeItem(pr.UpdatedAt)).
		withIntAttribute("PullRequestNumber", strconv.Itoa(pr.PullRequestNumber))

	itemAttrs := ad.allAttributes()

	return itemAttrs, nil
}

func (pr PullRequest) PartitionKey() string {
	return fmt.Sprintf("REPO#%s#%s", strings.ToLower(pr.RepoOwner), strings.ToLower(pr.RepoName))
}

func (pr PullRequest) initialAttributeDefinition() attrDefinition {
	sk := fmt.Sprintf("#PR#%s", pad(pr.PullRequestNumber))
	return attrDefinition{
		pk:        pr.PartitionKey(),
		sk:        sk,
		typeLabel: "PULLREQUEST",
	}
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

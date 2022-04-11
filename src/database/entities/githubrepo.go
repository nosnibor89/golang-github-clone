package entities

import (
	"fmt"
	"github-clone/src/model"
	"github.com/aws/aws-sdk-go/aws"
	"strings"
	"time"
)

type GithubRepo struct {
	Entity
	Name, Owner, Description string
}

func NewGithubRepo(name, owner, description string) GithubRepo {
	return GithubRepo{
		Name:        name,
		Owner:       owner,
		Description: description,
		Entity: Entity{
			time.Now(),
			time.Now(),
		},
	}
}

func (repo GithubRepo) Key() Attrs {
	ad := repo.initialAttributeDefinition(repo.Owner, repo.Name)
	return ad.getPrimaryKey()
}

// ToItem Exports entity to Item type
func (repo GithubRepo) ToItem() (Attrs, error) {
	//TODO: ADD VALIDATION
	gs1 := fmt.Sprintf("REPO#%s", repo.Name)
	ad := repo.initialAttributeDefinition(repo.Owner, repo.Name)
	ad.withStringAttribute("Name", repo.Name).
		withStringAttribute("Owner", repo.Owner).
		withStringAttribute("Description", repo.Description).
		withIntAttribute("IssuePRNumber", "0").
		withStringAttribute("CreatedAt", parseTimeItem(repo.CreatedAt)).
		withStringAttribute("UpdatedAt", parseTimeItem(repo.UpdatedAt)).
		withSecondaryIndexKey(1, gs1, gs1)

	itemAttrs := ad.allAttributes()

	return itemAttrs, nil
}

func (repo GithubRepo) ToModelFromAttrs(attrs Attrs) model.Repo {
	var repoModel model.Repo
	if attrs["Name"] != nil && attrs["Owner"] != nil {
		owner := aws.StringValue(attrs["Owner"].S)
		repoModel = model.Repo{
			Model: model.Model{
				Identifier: aws.StringValue(attrs["Name"].S),
			},
			Name:        aws.StringValue(attrs["Name"].S),
			Owner:       model.User{Username: owner, Name: owner},
			Description: aws.StringValue(attrs["Description"].S),
			UpdatedAt:   parseTimeAttr(aws.StringValue(attrs["UpdatedAt"].S)),
			CreatedAt:   parseTimeAttr(aws.StringValue(attrs["CreatedAt"].S)),
		}

	}
	return repoModel
}

func (repo GithubRepo) ToModel() model.Repo {
	repoModel := model.Repo{
		Model: model.Model{
			Identifier: repo.Name,
		},
		Name:        repo.Name,
		Owner:       model.User{Username: repo.Owner, Name: repo.Owner},
		Description: repo.Description,
		UpdatedAt:   repo.UpdatedAt,
		CreatedAt:   repo.CreatedAt,
	}
	return repoModel
}

func (repo GithubRepo) initialAttributeDefinition(owner, name string) attrDefinition {
	pk := fmt.Sprintf("REPO#%s#%s", strings.ToLower(owner), strings.ToLower(name))
	return attrDefinition{
		pk:        pk,
		sk:        pk,
		typeLabel: "REPO",
	}
}

package entities

import (
	"fmt"
	"github-clone/src/model"
	util2 "github-clone/src/util"
	"github.com/aws/aws-sdk-go/aws"
	"strings"
	"time"
)

const (
	itemType = "REPO"
)

type GithubRepo struct {
	Entity
	Name        string
	Owner       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (repo GithubRepo) Key() Attrs {
	pk := fmt.Sprintf("REPO#%s#%s", strings.ToLower(repo.Owner), strings.ToLower(repo.Name))
	ad := attrDefinition{
		pk:        pk,
		sk:        pk,
		typeLabel: itemType,
	}

	return ad.getPrimaryKey()
}

// ToItem Exports entity to Item type
func (repo GithubRepo) ToItem() (Attrs, error) {
	//TODO: ADD VALIDATION
	pk := fmt.Sprintf("REPO#%s#%s", strings.ToLower(repo.Owner), strings.ToLower(repo.Name))
	gs1 := fmt.Sprintf("REPO#%s", repo.Name)
	ad := attrDefinition{
		pk:        pk,
		sk:        pk,
		typeLabel: itemType,
	}
	ad.withStringAttribute("Name", repo.Name).
		withStringAttribute("Owner", repo.Owner).
		withStringAttribute("Description", repo.Description).
		withStringAttribute("CreatedAt", util2.ParseTimeItem(repo.CreatedAt)).
		withStringAttribute("UpdatedAt", util2.ParseTimeItem(repo.UpdatedAt)).
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
			UpdatedAt:   util2.ParseTimeAttr(aws.StringValue(attrs["UpdatedAt"].S)),
			CreatedAt:   util2.ParseTimeAttr(aws.StringValue(attrs["CreatedAt"].S)),
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

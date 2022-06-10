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

func (r GithubRepo) Key() Attrs {
	ad := r.initialAttributeDefinition(r.Owner, r.Name)
	return ad.getPrimaryKey()
}

// ToItem Exports entity to Item type
func (r GithubRepo) ToItem() (Attrs, error) {
	//TODO: ADD VALIDATION
	gs1 := fmt.Sprintf("REPO#%s#%s", strings.ToLower(r.Owner), strings.ToLower(r.Name))
	ad := r.initialAttributeDefinition(r.Owner, r.Name)
	ad.withStringAttribute("Name", r.Name).
		withStringAttribute("Owner", r.Owner).
		withStringAttribute("Description", r.Description).
		withIntAttribute("IssuePRNumber", "0").
		withStringAttribute("CreatedAt", parseTimeItem(r.CreatedAt)).
		withStringAttribute("UpdatedAt", parseTimeItem(r.UpdatedAt)).
		withSecondaryIndexKey(1, gs1, gs1)

	itemAttrs := ad.allAttributes()

	return itemAttrs, nil
}

func (r GithubRepo) ToModelFromAttrs(attrs Attrs) model.Repo {
	var repo model.Repo
	if attrs["Name"] != nil && attrs["Owner"] != nil {
		owner := aws.StringValue(attrs["Owner"].S)
		repo = model.Repo{
			Model: model.Model{
				UpdatedAt: parseTimeAttr(aws.StringValue(attrs["UpdatedAt"].S)),
				CreatedAt: parseTimeAttr(aws.StringValue(attrs["CreatedAt"].S)),
			},
			Name:        aws.StringValue(attrs["Name"].S),
			Owner:       model.User{Username: owner, Name: owner},
			Description: aws.StringValue(attrs["Description"].S),
		}

	}
	return repo
}

func (r GithubRepo) ToModel() model.Repo {
	repo := model.Repo{
		Model: model.Model{
			UpdatedAt: r.UpdatedAt,
			CreatedAt: r.CreatedAt,
		},
		Name:        r.Name,
		Owner:       model.User{Username: r.Owner, Name: r.Owner},
		Description: r.Description,
	}
	return repo
}

func (r GithubRepo) initialAttributeDefinition(owner, name string) attrDefinition {
	pk := fmt.Sprintf("REPO#%s#%s", strings.ToLower(owner), strings.ToLower(name))
	return attrDefinition{
		pk:        pk,
		sk:        pk,
		typeLabel: "REPO",
	}
}

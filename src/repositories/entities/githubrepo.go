package entities

import (
	"fmt"
	"github-clone/src/model"
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
		withStringAttribute("CreatedAt", repo.CreatedAt.String()).
		withStringAttribute("UpdatedAt", repo.UpdatedAt.String()).
		withSecondaryIndexKey(1, gs1, gs1)

	itemAttrs := ad.allAttributes()

	fmt.Println(itemAttrs)

	return itemAttrs, nil
}

func (repo GithubRepo) ToModel(_ Attrs) model.Repo {
	repoModel := model.Repo{
		Name:        repo.Name,
		Owner:       model.User{Username: repo.Owner, Name: repo.Owner},
		Description: repo.Description,
		UpdatedAt:   repo.UpdatedAt,
		CreatedAt:   repo.CreatedAt,
	}
	return repoModel
}

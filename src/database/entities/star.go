package entities

import (
	"fmt"
	"strings"
	"time"
)

type Star struct {
	Entity
	RepoName, RepoOwner, Username string
}

func NewStar(repoName, repoOwner, username string) Star {
	return Star{
		RepoName:  repoName,
		RepoOwner: repoOwner,
		Username:  username,
		Entity: Entity{
			time.Now(),
			time.Now(),
		},
	}
}

func (s Star) PartitionKey() string {
	return fmt.Sprintf("REPO#%s#%s", strings.ToLower(s.RepoOwner), strings.ToLower(s.RepoName))
}

// ToItem Exports entity to Item type
func (s Star) ToItem() (Attrs, error) {
	//TODO: ADD VALIDATION
	ad := s.initialAttributeDefinition()
	ad.withStringAttribute("RepoOwner", s.RepoOwner).
		withStringAttribute("RepoName", s.RepoName).
		withStringAttribute("Username", s.RepoOwner).
		withStringAttribute("CreatedAt", parseTimeItem(s.CreatedAt)).
		withStringAttribute("UpdatedAt", parseTimeItem(s.UpdatedAt))

	itemAttrs := ad.allAttributes()

	return itemAttrs, nil
}

func (s Star) initialAttributeDefinition() attrDefinition {
	sk := fmt.Sprintf("STAR#%s", strings.ToLower(s.Username))
	return attrDefinition{
		pk:        s.PartitionKey(),
		sk:        sk,
		typeLabel: "STAR",
	}
}

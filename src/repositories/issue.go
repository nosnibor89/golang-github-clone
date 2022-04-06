package repositories

import (
	"github-clone/src/model"
	"log"
)

type IssueRepository struct {
}

func (repo IssueRepository) Create(_ model.Issue) (model.Issue, error) {
	log.Printf("creating issue")
	return model.Issue{}, nil
}

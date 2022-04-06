package model

import "time"

type Issue struct {
	Model
	IssueNumber string    `json:"issueNumber"`
	Creator     User      `json:"creator"`
	Repo        Repo      `json:"repo"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (issue *Issue) WithCreator(user User) {
	issue.Creator = user
}

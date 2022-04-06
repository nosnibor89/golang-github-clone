package model

import "time"

type Repo struct {
	Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Owner       User      `json:"owner"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

package model

import "time"

type Repo struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Owner       User      `json:"owner"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

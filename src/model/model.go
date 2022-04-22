package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"time"
)

var (
	EncodingError = errors.New("could encode content")
	DecodingError = errors.New("could decode content")
)

type Model struct {
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type RepoElement struct {
	Creator User   `json:"creator"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Open    bool   `json:"open"`
	Repo    Repo   `json:"-"`
}

type Issue struct {
	Model
	RepoElement
	IssueNumber int `json:"issueNumber"`

	//TODO: Delete if everything is working
	//Creator     User   `json:"creator"`
	//Title       string `json:"title"`
	//Content     string `json:"content"`
	//Open        bool   `json:"open"`
	//Repo        Repo   `json:"-"`
}

func (re *RepoElement) WithCreator(user User) {
	re.Creator = user
}

func (i *Issue) FromJSON(json string) error {
	return parseToModel(json, i)
}

func (i Issue) ToJSON() (string, error) {
	return parseToJson(i)
}

type Repo struct {
	Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       User   `json:"owner"`
}

func (r *Repo) FromJSON(json string) error {
	return parseToModel(json, r)
}

func (r Repo) ToJSON() (string, error) {
	return parseToJson(r)
}

type PullRequest struct {
	Model
	RepoElement
	PullRequestNumber int `json:"pullRequestNumber"`

	//TODO: Delete if everything is working
	//Creator     User   `json:"creator"`
	//Title       string `json:"title"`
	//Content     string `json:"content"`
	//Open        bool   `json:"open"`
	//Repo        Repo   `json:"-"`
}

func (pr *PullRequest) FromJSON(json string) error {
	return parseToModel(json, pr)
}

func (pr PullRequest) ToJSON() (string, error) {
	return parseToJson(pr)
}

type User struct {
	Model
	Name     string `json:"name"`
	Username string `json:"username"`
}

func GetUserFromRequest(request events.APIGatewayProxyRequest) User {
	user := request.Headers["Authorization"]
	return User{
		Username: user,
		Name:     user,
	}
}

func parseToModel(content string, modelValue interface{}) error {
	if err := json.Unmarshal([]byte(content), &modelValue); err != nil {
		msg := fmt.Sprintf("Could not parse body correctly %v\n", err)
		return fmt.Errorf("%s %w", msg, DecodingError)
	}

	return nil
}

func parseToJson(modelValue interface{}) (string, error) {
	decoded, err := json.Marshal(modelValue)
	if err != nil {
		msg := fmt.Sprintf("Could not parse created repo %v\n", err)
		return "", fmt.Errorf("%s %w", msg, EncodingError)
	}

	return string(decoded), nil
}

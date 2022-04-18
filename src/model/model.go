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

type Issue struct {
	Model
	IssueNumber int    `json:"issueNumber"`
	Creator     User   `json:"creator"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Open        bool   `json:"open"`
	Repo        Repo   `json:"-"`
}

func (issue *Issue) WithCreator(user User) {
	issue.Creator = user
}

func (issue *Issue) FromJSON(json string) error {
	return parseToModel(json, issue)
}

func (issue Issue) ToJSON() (string, error) {
	return parseToJson(issue)
}

type Repo struct {
	Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       User   `json:"owner"`
}

func (model *Repo) FromJSON(json string) error {
	return parseToModel(json, model)
}

func (model Repo) ToJSON() (string, error) {
	return parseToJson(model)
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

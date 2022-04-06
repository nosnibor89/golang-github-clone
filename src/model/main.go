package model

import (
	"encoding/json"
	"fmt"
	"github-clone/src/util"
	"reflect"
)

//TODO: DO WE NEED THIS ??
// TypeConstraint maybe is not the best name
//type TypeConstraint interface {
//	Repo | Issue | User
//}

type Model struct {
	Identifier string
}

type DecodingError struct {
	Error   error
	Code    int
	Message string
}

func (model Model) FromJSON(json string) DecodingError {
	if err, code, msg := parseToModel(json, &model); err != nil {
		return DecodingError{
			Error:   err,
			Code:    code,
			Message: msg,
		}
	}

	return DecodingError{nil, 0, ""}
}

func (model Model) ToJSON() (DecodingError, string) {
	err, code, jsonContent := parseToJson(model)

	if err != nil {
		return DecodingError{
			Error:   err,
			Code:    code,
			Message: jsonContent,
		}, ""
	}

	return DecodingError{nil, 0, ""}, jsonContent
}

func parseToModel(content string, modelValue *Model) (error, int, string) {
	if err := json.Unmarshal([]byte(content), &modelValue); err != nil {
		msg := fmt.Sprintf("Could not parse body correctly %v\n", err)
		fmt.Println(msg)
		statusCode := util.HttpErrorFromException(err).Code

		return err, statusCode, msg
	}

	return nil, 0, ""
}

func parseToJson(modelValue Model) (error, int, string) {
	decoded, err := json.Marshal(modelValue)
	if err != nil {
		msg := fmt.Sprintf("Could not parse created repo %v\n", err)
		fmt.Println(msg)
		statusCode := util.HttpErrorFromException(err).Code

		return err, statusCode, msg
	}

	return nil, 0, string(decoded)
}

func (model Model) IsEmpty() bool {
	return reflect.ValueOf(model).IsZero()
}

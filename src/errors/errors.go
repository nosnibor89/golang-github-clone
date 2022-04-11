package errors

import (
	"errors"
	"github-clone/src/model"
	"strings"
)

type HttpError struct {
	Code    int
	Message string
	Err     error
}

func (m HttpError) Error() string {
	return m.Message
}

func HttpErrorFromException(e error) HttpError {
	if strings.Contains(e.Error(), "ConditionalCheckFailedException:") {
		return HttpError{
			Code:    409,
			Message: "entity already exists",
		}
	}

	if errors.Is(e, model.EncodingError) {
		return HttpError{
			Code:    400,
			Message: "could parse content in request correctly",
		}
	}

	return HttpError{
		Code:    500,
		Message: "error creating entity",
	}
}

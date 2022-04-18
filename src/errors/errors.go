package errors

import (
	"errors"
	"github-clone/src/model"
	"net/http"
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
			Code:    http.StatusConflict,
			Message: "item already exists",
		}
	}

	if errors.Is(e, model.EncodingError) {
		return HttpError{
			Code:    http.StatusBadRequest,
			Message: "could parse content in request correctly",
		}
	}

	return HttpError{
		Code:    http.StatusInternalServerError,
		Message: "error processing request",
	}
}

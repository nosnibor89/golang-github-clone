package errors

import (
	"errors"
	"github-clone/src/model"
	"log"
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
	log.Printf("ERROR: %v", e)
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

	if strings.Contains(e.Error(), "ValidationException:") {
		return HttpError{
			Code:    http.StatusBadRequest,
			Message: "incorrect request",
		}
	}

	return HttpError{
		Code:    http.StatusInternalServerError,
		Message: "error processing request",
	}
}

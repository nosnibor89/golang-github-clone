package errors

import (
	"errors"
	"fmt"
	"github-clone/src/model"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

	switch t := errors.Unwrap(e).(type) {
	case *dynamodb.ConditionalCheckFailedException:
	case *dynamodb.TransactionCanceledException:
		fmt.Println(t.CancellationReasons)
		return HttpError{
			Code:    http.StatusConflict,
			Message: "item already exists",
		}
	}

	switch t := e.(type) {
	case *dynamodb.ConditionalCheckFailedException:
	case *dynamodb.TransactionCanceledException:
		fmt.Println(t.CancellationReasons)
		return HttpError{
			Code:    http.StatusConflict,
			Message: "item already exists",
		}
	}

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

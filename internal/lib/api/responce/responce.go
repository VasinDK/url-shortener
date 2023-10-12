package responce

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

type Responce struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

func Ok() Responce {
	return Responce{
		Status: StatusOk,
	}
}

func Error(msg string) Responce {
	return Responce{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Responce {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid url", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Responce{
		Status: StatusError,
		Error:  strings.Join(errMsgs, " ,"),
	}
}

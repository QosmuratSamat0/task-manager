package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "ERROR"
)

func OK() Response {
	return Response{Status: StatusOK}
}

func Error(msg string) Response {
	return Response{Status: StatusError, Error: msg}

}

func ValidationError(errs validator.ValidationErrors) Response {
	var errsMsg []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errsMsg = append(errsMsg, fmt.Sprintf("%s is a required field", err.Field()))
		case "user":
			errsMsg = append(errsMsg, fmt.Sprintf("%s is not a valid user", err.Field()))
		case "task":
			errsMsg = append(errsMsg, fmt.Sprintf("%s is not a valid task", err.Field()))
		default:
			errsMsg = append(errsMsg, fmt.Sprintf("%s is not a valid", err.Field()))
		}

	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errsMsg, ","),
	}
}

package response

import (
	"fmt"
	"github.com/go-playground/validator"
	"strings"
)

type Response struct {
	Error string `json:"error,omitempty"`
}

func ErrorResponse(msg string) *Response {
	return &Response{
		Error: msg,
	}
}

func OkResponse() *Response {
	return &Response{}
}

func ValidationError(errs validator.ValidationErrors) *Response {
	var errList []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errList = append(errList, fmt.Sprintf("field %s is required", err.Field()))
		case "url":
			errList = append(errList, fmt.Sprintf("field %s is not a valid url", err.Field()))
		default:
			errList = append(errList, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return &Response{
		Error: strings.Join(errList, ", "),
	}
}

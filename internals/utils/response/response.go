package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	StatusOk = "OK"
	StatusErrr = "Error"
)

type Response struct {
	Status  string      `json:"status"`
	Error   string        `json:"error,omitempty"`
}

func WriteJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusOk,
		Error:  err.Error(),
	}
}

func ValidateError(errs validator.ValidationErrors) Response {
	var errsString []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errsString = append(errsString, fmt.Sprintf("field %s is required", err.Field()))
		default:
			errsString = append(errsString, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return Response{
		Status: StatusErrr,
		Error: strings.Join(errsString, ", "),
	}
}
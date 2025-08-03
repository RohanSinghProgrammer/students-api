package response

import (
	"encoding/json"
	"net/http"
)

const (
	StatusOk = "OK"
	StatusErrr = "Error"
)

type Response struct {
	Status  string      `json:"status"`
	Error   error       `json:"error,omitempty"`
}

func WriteJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusOk,
		Error:  err,
	}
}
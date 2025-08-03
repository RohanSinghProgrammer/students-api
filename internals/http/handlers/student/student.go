package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/rohansinghprogrammer/sudents-api/internals/types"
	"github.com/rohansinghprogrammer/sudents-api/internals/utils/response"
)

func New () http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err = response.WriteJson(w, http.StatusCreated, map[string]string{
			"message": "Student created successfully",
			"id":      student.ID,
		})

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}

		slog.Info("Student created successfully", slog.String("id", student.ID))
	}
} 
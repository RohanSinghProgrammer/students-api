package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/rohansinghprogrammer/sudents-api/internals/storage"
	"github.com/rohansinghprogrammer/sudents-api/internals/types"
	"github.com/rohansinghprogrammer/sudents-api/internals/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student
		validate := validator.New()

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// validate request body
		if err := validate.Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidateError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{
			"Student Created Successfuly! , ID: ": int64(lastId),
		})

		slog.Info("Student created successfully", slog.Int64("lastId", int64(lastId)))
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			slog.Error("Invalid ID")
			response.WriteJson(w, http.StatusBadRequest, map[string]string{"Invalid ID": id})
			return
		}
		student, err := storage.GetStudentById(parsedId)
		if err != nil {
			slog.Error("Student doesn't exists")
			response.WriteJson(w, http.StatusInternalServerError, map[string]string{"Error": "Student not found"})
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := storage.GetStudentsList()
		if err != nil {
			slog.Error("Error getting students")
			response.WriteJson(w, http.StatusInternalServerError, err.Error())
			return
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}

func Update(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			slog.Error("Invalid ID")
			response.WriteJson(w, http.StatusBadRequest, map[string]string{"Invalid ID": id})
			return
		}

		var student types.Student
		validate := validator.New()

		err = json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validate.Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidateError(validateErrs))
			return
		}

		updatedStudent, err := storage.UpdateStudent(
			parsedId,
			student.Name,
			student.Email,
			student.Age,
		)

		if err != nil {
			slog.Error("Error updating student", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, updatedStudent)
		slog.Info("Student updated successfully", slog.Uint64("id", parsedId))
	}
}

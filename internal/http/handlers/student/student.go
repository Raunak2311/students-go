package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/raunak173/students-go/internal/storage"
	"github.com/raunak173/students-go/internal/types"
	"github.com/raunak173/students-go/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Creating a student")

		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//Request validation
		err = validator.New().Struct(student)
		if err != nil {
			validateErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErr))
			return
		}

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id")
		slog.Info("Getting a student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Getting all students")

		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}

func UpdateById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id")
		slog.Info("Getting a student", slog.String("id", id))

		// Parse the request body into a Student object
		var student types.Student
		if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
			resp := response.GeneralError(fmt.Errorf("invalid request body"))
			response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		// Validate the student data if necessary
		v := validator.New()
		if err := v.Struct(student); err != nil {
			if validationErrs, ok := err.(validator.ValidationErrors); ok {
				resp := response.ValidationError(validationErrs)
				response.WriteJson(w, http.StatusBadRequest, resp)
				return
			}
			resp := response.GeneralError(err)
			response.WriteJson(w, http.StatusInternalServerError, resp)
			return
		}

		// Convert ID from string to int64 (assuming your ID is int64)
		studentId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			resp := response.GeneralError(fmt.Errorf("invalid id format"))
			response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		// Call the storage layer to update the student by ID
		updatedStudent, err := storage.UpdateStudentById(studentId, student)
		if err != nil {
			resp := response.GeneralError(fmt.Errorf("failed to update student: %w", err))
			response.WriteJson(w, http.StatusInternalServerError, resp)
			return
		}
		response.WriteJson(w, http.StatusOK, updatedStudent)
	}
}

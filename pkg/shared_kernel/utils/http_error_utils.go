package utils

import (
	"elex_storage/pkg/shared_kernel/models"
	"encoding/json"
	"errors"
	"net/http"
)

type httpError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type HttpErrorUtils struct {
}

func NewHttpErrorUtils() *HttpErrorUtils {
	return &HttpErrorUtils{}
}

// NotFound sends a 404 Not Found response with a custom message.
func (h *HttpErrorUtils) NotFound(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(httpError{Message: message, Code: http.StatusNotFound})
}

// BadRequest sends a 400 Bad Request response with a custom message.
func (h *HttpErrorUtils) BadRequest(w http.ResponseWriter, err error) {
	var commonErr *models.CommonError
	if !errors.As(err, &commonErr) {
		err := models.NewCommonError(nil)
		err.AddErr("_", err.Error())
	}
	w.Write([]byte(err.Error()))
}

// InternalServerError sends a 500 Internal Server Error response with a custom message.
func (h *HttpErrorUtils) InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(httpError{Message: "The server has encountered a situation it does not know how to handle. This error is generic, indicating that the server cannot find a more appropriate 5XX status code to respond with.", Code: http.StatusInternalServerError})
}

func (h *HttpErrorUtils) StatusConflict(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(httpError{Message: message, Code: http.StatusInternalServerError})
}

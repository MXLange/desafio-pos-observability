package zipcode

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrNotFoundZipcode = errors.New(ErrNotFound)

type ErrorResponse struct {
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

func StatusFromError(err error) (int, string) {
	switch {
	case IsValidationError(err):
		return http.StatusUnprocessableEntity, ErrInvalidZipcode
	case errors.Is(err, ErrNotFoundZipcode):
		return http.StatusNotFound, ErrNotFound
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

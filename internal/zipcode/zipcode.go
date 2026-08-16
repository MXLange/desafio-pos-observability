package zipcode

import (
	"regexp"
	"strings"
)

const (
	ErrInvalidZipcode = "invalid zipcode"
	ErrNotFound       = "can not find zipcode"
)

var digitsOnly = regexp.MustCompile(`^\d{8}$`)

type Request struct {
	CEP string `json:"cep"`
}

func Validate(cep string) (string, error) {
	cep = strings.TrimSpace(cep)
	if !digitsOnly.MatchString(cep) {
		return "", errInvalid()
	}
	return cep, nil
}

func errInvalid() error {
	return &ValidationError{Message: ErrInvalidZipcode}
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

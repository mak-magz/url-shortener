package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("AppError: [%d] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func NewNotFoundError(message string, err error) *AppError {
	return NewAppError(http.StatusNotFound, message, err)
}

func NewInternalError(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, message, err)
}

func NewBadRequestError(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, message, err)
}

func FormatValidationError(err error) *AppError {
	var errs []string

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			errs = append(errs, e.Field()+" is invalid or missing.")
		}

		return NewAppError(http.StatusBadRequest, "Validation Error", errors.New(strings.Join(errs, ", ")))
	}

	return NewAppError(http.StatusBadRequest, "Validation Error", errors.New("Invalid request data"))
}

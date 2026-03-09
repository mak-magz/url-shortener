package errors_test

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	appErrors "github.com/mak-magz/url-shortener/platform/errors"
)

func TestNewAppError(t *testing.T) {
	err := errors.New("test error")
	appErr := appErrors.NewAppError(500, "test message", err)

	if appErr.Code != 500 {
		t.Errorf("Expected code 500, got %d", appErr.Code)
	}

	if appErr.Message != "test message" {
		t.Errorf("Expected message 'test message', got %s", appErr.Message)
	}

	if appErr.Err != err {
		t.Errorf("Expected error 'test error', got %v", appErr.Err)
	}
}

func TestNewNotFoundError(t *testing.T) {
	err := errors.New("test error")
	appErr := appErrors.NewNotFoundError("test message", err)

	if appErr.Code != 404 {
		t.Errorf("Expected code 404, got %d", appErr.Code)
	}

	if appErr.Message != "test message" {
		t.Errorf("Expected message 'test message', got %s", appErr.Message)
	}

	if appErr.Err != err {
		t.Errorf("Expected error 'test error', got %v", appErr.Err)
	}
}

func TestNewInternalError(t *testing.T) {
	err := errors.New("test error")
	appErr := appErrors.NewInternalError("test message", err)

	if appErr.Code != 500 {
		t.Errorf("Expected code 500, got %d", appErr.Code)
	}

	if appErr.Message != "test message" {
		t.Errorf("Expected message 'test message', got %s", appErr.Message)
	}

	if appErr.Err != err {
		t.Errorf("Expected error 'test error', got %v", appErr.Err)
	}
}

func TestNewBadRequestError(t *testing.T) {
	err := errors.New("test error")
	appErr := appErrors.NewBadRequestError("test message", err)

	if appErr.Code != 400 {
		t.Errorf("Expected code 400, got %d", appErr.Code)
	}

	if appErr.Message != "test message" {
		t.Errorf("Expected message 'test message', got %s", appErr.Message)
	}

	if appErr.Err != err {
		t.Errorf("Expected error 'test error', got %v", appErr.Err)
	}
}

func TestAppErrorUnwrap(t *testing.T) {
	err := errors.New("test error")
	appErr := appErrors.NewAppError(500, "test message", err)

	if appErr.Unwrap() != err {
		t.Errorf("Expected error 'test error', got %v", appErr.Unwrap())
	}
}

func TestAppErrorError(t *testing.T) {
	err := errors.New("test error")
	appErr := appErrors.NewAppError(500, "test message", err)

	if appErr.Error() != "AppError: [500] test message" {
		t.Errorf("Expected error 'AppError: [500] test message', got %s", appErr.Error())
	}
}

func TestFormatValidationError(t *testing.T) {
	validate := validator.New()

	type TestStruct struct {
		Name string `validate:"required"`
	}

	t.Run("Validator Error", func(t *testing.T) {
		badStruct := TestStruct{}
		err := validate.Struct(badStruct)

		appErr := appErrors.FormatValidationError(err)

		if appErr.Code != 400 {
			t.Errorf("Expected code 400, got %d", appErr.Code)
		}

		if appErr.Message != "Validation Error" {
			t.Errorf("Expected message 'Validation Error', got %s", appErr.Message)
		}

		expectedDetails := "Name is invalid or missing."
		if appErr.Err == nil || appErr.Err.Error() != expectedDetails {
			t.Errorf("Expected error detail '%s', got %v", expectedDetails, appErr.Err)
		}
	})

	t.Run("Generic Error", func(t *testing.T) {
		genericErr := errors.New("some generic error context")

		appErr := appErrors.FormatValidationError(genericErr)

		if appErr.Code != 400 {
			t.Errorf("Expected code 400, got %d", appErr.Code)
		}

		if appErr.Message != "Validation Error" {
			t.Errorf("Expected message 'Validation Error', got %s", appErr.Message)
		}

		expectedDetails := "Invalid request data"
		if appErr.Err == nil || appErr.Err.Error() != expectedDetails {
			t.Errorf("Expected error '%s', got %v", expectedDetails, appErr.Err)
		}
	})
}

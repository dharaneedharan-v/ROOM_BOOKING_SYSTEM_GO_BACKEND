package utils_test
import (
	"BookingSystem/Booking/internal/dtos"
	errorcodes "BookingSystem/Booking/internal/errorcodes"
	"BookingSystem/Booking/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertHeader(t *testing.T, h http.Header, key, expected string) {
	if got := h.Get(key); got != expected {
		t.Errorf("handler returned wrong %s header: got %v want %v", key, got, expected)
	}
}

func TestWriteResponse(t *testing.T) {
	t.Run("Success response", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"message": "success"}

		utils.WriteResponse(w, http.StatusOK, data)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"message":"success"}`, w.Body.String())
	})

	t.Run("Database Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"message": "internal server"}

		utils.WriteResponse(w, http.StatusInternalServerError, data)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"message":"internal server"}`, w.Body.String())
	})

	t.Run("Item Not Found Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"message": "data not found"}

		utils.WriteResponse(w, http.StatusNotFound, data)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"message":"data not found"}`, w.Body.String())
	})

	t.Run("Error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"error": "bad request"}

		utils.WriteResponse(w, http.StatusBadRequest, data)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"error":"bad request"}`, w.Body.String())
	})

	t.Run("Json Encode Error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := make(chan int)

		utils.WriteResponse(w, http.StatusInternalServerError, data)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestValidateStruct(t *testing.T) {
	type TestStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	t.Run("Valid struct", func(t *testing.T) {
		input := TestStruct{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		errors := utils.ValidateStruct(&input)
		assert.Nil(t, errors)
	})

	t.Run("Missing name", func(t *testing.T) {
		input := TestStruct{
			Email: "john@example.com",
		}

		expectedErrors := []dtos.Error{
			{
				Field:   "Name",
				Message: "Validation failed for the field 'Name' on the 'required' constraint.",
				Code:    errorcodes.ErrorCodeStatus[errorcodes.ValidationErrorCode],
			},
		}

		errors := utils.ValidateStruct(&input)
		assert.Equal(t, expectedErrors, errors)
	})

	t.Run("Invalid email", func(t *testing.T) {
		input := TestStruct{
			Name:  "John Doe",
			Email: "invalid-email",
		}

		expectedErrors := []dtos.Error{
			{
				Field:   "Email",
				Message: "Validation failed for the field 'Email' on the 'email' constraint.",
				Code:    errorcodes.ErrorCodeStatus[errorcodes.ValidationErrorCode],
			},
		}

		errors := utils.ValidateStruct(&input)
		assert.Equal(t, expectedErrors, errors)
	})

	t.Run("Empty struct", func(t *testing.T) {
		input := TestStruct{}

		expectedErrors := []dtos.Error{
			{
				Field:   "Name",
				Message: "Validation failed for the field 'Name' on the 'required' constraint.",
				Code:    errorcodes.ErrorCodeStatus[errorcodes.ValidationErrorCode],
			},
			{
				Field:   "Email",
				Message: "Validation failed for the field 'Email' on the 'required' constraint.",
				Code:    errorcodes.ErrorCodeStatus[errorcodes.ValidationErrorCode],
			},
		}

		errors := utils.ValidateStruct(&input)
		assert.ElementsMatch(t, expectedErrors, errors)
	})
}

func TestMapErrorCode(t *testing.T) {
	t.Run("Record Not Found Error", func(t *testing.T) {
		input := &dtos.APIResponse{
			Errors: []dtos.Error{
				{
					Code: errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode],
				},
			},
		}

		expectedStatus := http.StatusNotFound
		result := utils.MapErrorCode(input)
		assert.Equal(t, expectedStatus, result.Code)
	})

	t.Run("Validation Error", func(t *testing.T) {
		input := &dtos.APIResponse{
			Errors: []dtos.Error{
				{
					Code: errorcodes.ErrorCodeStatus[errorcodes.ValidationErrorCode],
				},
			},
		}

		expectedStatus := http.StatusBadRequest
		result := utils.MapErrorCode(input)
		assert.Equal(t, expectedStatus, result.Code)
	})

	t.Run("Unexpected Error", func(t *testing.T) {
		input := &dtos.APIResponse{
			Errors: []dtos.Error{
				{
					Code: errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
				},
			},
		}

		expectedStatus := http.StatusInternalServerError
		result := utils.MapErrorCode(input)
		assert.Equal(t, expectedStatus, result.Code)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		input := &dtos.APIResponse{
			Errors: []dtos.Error{
				{
					Code: errorcodes.ErrorCodeStatus[errorcodes.InternalServerErrorCode],
				},
			},
		}

		expectedStatus := http.StatusInternalServerError
		result := utils.MapErrorCode(input)
		assert.Equal(t, expectedStatus, result.Code)
	})

	t.Run("Database Not Up Error", func(t *testing.T) {
		input := &dtos.APIResponse{
			Errors: []dtos.Error{
				{
					Code: errorcodes.ErrorCodeStatus[errorcodes.DatabaseNotUpErrorCode],
				},
			},
		}

		expectedStatus := http.StatusInternalServerError
		result := utils.MapErrorCode(input)
		assert.Equal(t, expectedStatus, result.Code)
	})

	t.Run("Unauthorized Error", func(t *testing.T) {
		input := &dtos.APIResponse{
			Errors: []dtos.Error{
				{
					Code: errorcodes.ErrorCodeStatus[errorcodes.UnauthorizedErrorCode],
				},
			},
		}

		expectedStatus := http.StatusUnauthorized
		result := utils.MapErrorCode(input)
		assert.Equal(t, expectedStatus, result.Code)
	})

	t.Run("Missing Auth Header Error", func(t *testing.T) {
		input := &dtos.APIResponse{
			Errors: []dtos.Error{
				{
					Code: errorcodes.ErrorCodeStatus[errorcodes.MissingAuthHeaderErrorCode],
				},
			},
		}

		expectedStatus := http.StatusUnauthorized
		result := utils.MapErrorCode(input)
		assert.Equal(t, expectedStatus, result.Code)
	})

	t.Run("Invalid Auth Type Error", func(t *testing.T) {
		input := &dtos.APIResponse{
			Errors: []dtos.Error{
				{
					Code: errorcodes.ErrorCodeStatus[errorcodes.InvalidAuthTypeErrorCode],
				},
			},
		}

		expectedStatus := http.StatusUnauthorized
		result := utils.MapErrorCode(input)
		assert.Equal(t, expectedStatus, result.Code)
	})

	t.Run("Bad Request Error", func(t *testing.T) {
		input := &dtos.APIResponse{
			Errors: []dtos.Error{
				{
					Code: errorcodes.ErrorCodeStatus[errorcodes.BadRequestErrorCode],
				},
			},
		}

		expectedStatus := http.StatusBadRequest
		result := utils.MapErrorCode(input)
		assert.Equal(t, expectedStatus, result.Code)
	})

	t.Run("Unprocessable Entity Error", func(t *testing.T) {
		input := &dtos.APIResponse{
			Errors: []dtos.Error{
				{
					Code: errorcodes.ErrorCodeStatus[errorcodes.UnprocessableEntityErrorCode],
				},
			},
		}

		expectedStatus := http.StatusUnprocessableEntity
		result := utils.MapErrorCode(input)
		assert.Equal(t, expectedStatus, result.Code)
	})

	t.Run("Forbidden Error", func(t *testing.T) {
		input := &dtos.APIResponse{
			Errors: []dtos.Error{
				{
					Code: errorcodes.ErrorCodeStatus[errorcodes.ForbiddenErrorCode],
				},
			},
		}

		expectedStatus := http.StatusForbidden
		result := utils.MapErrorCode(input)
		assert.Equal(t, expectedStatus, result.Code)
	})

	t.Run("Unknown Error", func(t *testing.T) {
		input := &dtos.APIResponse{
			Errors: []dtos.Error{
				{
					Code: "UNKNOWN_ERROR",
				},
			},
		}

		expectedStatus := http.StatusInternalServerError
		result := utils.MapErrorCode(input)
		assert.Equal(t, expectedStatus, result.Code)
	})

	t.Run("No Error", func(t *testing.T) {
		input := &dtos.APIResponse{
			Errors: []dtos.Error{},
		}

		expectedStatus := http.StatusOK
		result := utils.MapErrorCode(input)
		assert.Equal(t, expectedStatus, result.Code)
	})
}



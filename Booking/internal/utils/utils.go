
package utils

import (
	"encoding/json"
	"fmt"
	"BookingSystem/Booking/internal/dtos"
	"BookingSystem/Booking/internal/errorcodes"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func WriteResponse(w http.ResponseWriter, code int, data interface{}) {
	// Don't overwrite existing headers
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}
	// Don't set CORS headers here, they should be set by the EnableCors middleware

	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validate = validator.New()

func ValidateStruct(s interface{}) []dtos.Error {
	var validationErrors []dtos.Error
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, dtos.Error{
				Field:   err.Field(),
				Message: fmt.Sprintf("Validation failed for the field '%s' on the '%s' constraint.", err.Field(), err.Tag()),
				Code:    errorcodes.ErrorCodeStatus[errorcodes.ValidationErrorCode],
			})
		}
	}
	return validationErrors
}

func MapErrorCode(apiResponse *dtos.APIResponse) *dtos.APIResponse {
	if len(apiResponse.Errors) == 0 {
		apiResponse.Code = http.StatusOK
		return apiResponse
	}

	switch apiResponse.Errors[0].Code {
	case errorcodes.ErrorCodeStatus[errorcodes.RecordNotFoundErrorCode]:
		apiResponse.Code = http.StatusNotFound
		return apiResponse
	case errorcodes.ErrorCodeStatus[errorcodes.ValidationErrorCode]:
		apiResponse.Code = http.StatusBadRequest
		return apiResponse
	case errorcodes.ErrorCodeStatus[errorcodes.UnexpectedErrorCode],
		 errorcodes.ErrorCodeStatus[errorcodes.InternalServerErrorCode],
		 errorcodes.ErrorCodeStatus[errorcodes.DatabaseNotUpErrorCode]:
		apiResponse.Code = http.StatusInternalServerError
		return apiResponse
	case errorcodes.ErrorCodeStatus[errorcodes.UnauthorizedErrorCode],
		 errorcodes.ErrorCodeStatus[errorcodes.MissingAuthHeaderErrorCode],
		 errorcodes.ErrorCodeStatus[errorcodes.InvalidAuthTypeErrorCode]:
		apiResponse.Code = http.StatusUnauthorized
		return apiResponse
	case errorcodes.ErrorCodeStatus[errorcodes.BadRequestErrorCode]:
		apiResponse.Code = http.StatusBadRequest
		return apiResponse
	case errorcodes.ErrorCodeStatus[errorcodes.UnprocessableEntityErrorCode]:
		apiResponse.Code = http.StatusUnprocessableEntity
		return apiResponse
	case errorcodes.ErrorCodeStatus[errorcodes.ForbiddenErrorCode]:
		apiResponse.Code = http.StatusForbidden
		return apiResponse
	default:
		apiResponse.Code = http.StatusInternalServerError
		return apiResponse
	}
}

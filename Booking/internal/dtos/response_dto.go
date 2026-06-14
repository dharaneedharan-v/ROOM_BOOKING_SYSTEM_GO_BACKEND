package dtos

// APIResponse is the standardized response structure for all API endpoints
type APIResponse struct {
	Status    string      `json:"status"`           // "success" or "error"
	Code      int         `json:"code"`             // HTTP status code
	Message   string      `json:"message"`          // Descriptive message
	Data      interface{} `json:"data,omitempty"`   // Actual data (optional, hence 'omitempty')
	Errors    []Error     `json:"errors,omitempty"` // Array of errors (optional)
	Meta      *Meta       `json:"meta,omitempty"`   // Pagination or additional metadata (optional)
	RequestID string      `json:"request_id"`       // Unique Request ID for tracking requests
	Timestamp string      `json:"timestamp"`        // Timestamp
}


// Error represents a specific error in the API response
type Error struct {
	Field   string `json:"field,omitempty"` // Field where the error occurred
	Message string `json:"message"`         // Descriptive error message
	Code    string `json:"code"`            // Custom error code
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.Message
}

// Meta contains pagination information
type Meta struct {
	Page     int `json:"page"`      // Current page number
	PageSize int `json:"page_size"` // Number of items per page
	Total    int `json:"total"`     // Total number of items
}

// ValidationError structure for handling validation errors
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// ValidationErrorToAPIError converts a ValidationError to an API Error
func ValidationErrorToAPIError(err *ValidationError) Error {
	return Error{
		Field:   err.Field,
		Message: err.Message,
		Code:    "validation_error",
	}
}

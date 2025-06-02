package response

import (
	"net/http"
)

// APIResponse represents a standardized API response.
type APIResponse struct {
	Code    int         `json:"code"`    // Business-specific status code
	Message string      `json:"message"` // Message associated with the code
	Data    interface{} `json:"data"`    // Response data
}

// Success creates a new successful APIResponse.
// Uses http.StatusOK (200) for the HTTP status code and a business code of 0 for success.
func Success(data interface{}) APIResponse {
	return APIResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

// Fail creates a new failed APIResponse.
// The HTTP status code will be determined by the ErrorHandler middleware based on the error.
func Fail(code int, msg string) APIResponse {
	return APIResponse{
		Code:    code,
		Message: msg,
		Data:    nil,
	}
}

// Result is a generic helper to return JSON response
func Result(code int, data interface{}, msg string, c HTTPContext) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    code,
		Data:    data,
		Message: msg,
	})
}

// HTTPContext provides an abstraction over gin.Context for JSON responses.
type HTTPContext interface {
	JSON(code int, obj interface{})
}

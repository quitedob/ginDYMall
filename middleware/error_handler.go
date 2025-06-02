package middleware

import (
	"douyin/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler creates a middleware to handle errors globally.
func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next() // Process request

		// Check for errors
		errs := ctx.Errors
		if len(errs) > 0 {
			err := errs.Last() // Get the last error

			// TODO: Consider logging the error here

			// Default error code and message
			// We can define more specific error codes and messages based on error types
			// For now, using a generic error code for unhandled errors
			apiResp := response.Fail(http.StatusInternalServerError, err.Error())

			// Check if the error is a custom error type or if we can infer a status code
			// For example, if using custom error structs that include a status code:
			// if customErr, ok := err.Err.(MyCustomError); ok {
			//    apiResp.Code = customErr.BusinessCode()
			//    ctx.JSON(customErr.StatusCode(), apiResp)
			//    return
			// }

			// If it's a validation error or similar, it might have already been handled
			// and ctx.Writer.Written() would be true.
			// However, we are aiming to standardize all error responses.
			// For now, we assume if ctx.Errors has items, we need to send our standard response.

			// Respond with JSON
			// Note: If a response has already been written, this might not work as expected
			// or might cause "multiple response.WriteHeader calls" errors.
			// Gin's behavior is to ignore subsequent writes if the header has been written.
			if !ctx.Writer.Written() {
				ctx.JSON(http.StatusInternalServerError, apiResp) // Default to 500 for now
			}
			return // Stop further processing if error handled
		}
	}
}

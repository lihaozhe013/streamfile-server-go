package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse is the error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse is the success response structure
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SendError sends an error response
func SendError(c *gin.Context, code int, message string, details ...string) {
	response := ErrorResponse{
		Error: message,
	}

	if len(details) > 0 && details[0] != "" {
		response.Details = details[0]
	}

	c.JSON(code, response)
}

// SendSuccess sends a success response
func SendSuccess(c *gin.Context, message string, data interface{}) {
	response := SuccessResponse{
		Message: message,
		Data:    data,
	}

	c.JSON(http.StatusOK, response)
}

// SendJSON sends a JSON response
func SendJSON(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

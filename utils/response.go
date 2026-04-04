package utils

import (
	"github.com/gin-gonic/gin"
)

// Response standar untuk semua API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func SendError(c *gin.Context, statusCode int, message string, err interface{}) {
	var detail interface{}

	switch v := err.(type) {
	case error:
		detail = v.Error()
	default:
		detail = v
	}

	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   detail,
	})
}

func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

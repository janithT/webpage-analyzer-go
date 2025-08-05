package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteError sends a standardized error response to the client
func WriteError(ginC *gin.Context, statusCode int, message string) {
	ginC.JSON(statusCode, BaseResponse{
		Status:  "error",
		Message: message,
	})
	ginC.Abort()
}

// WriteSuccess sends a standardized success response with data
func WriteSuccess(ginC *gin.Context, message string, data interface{}) {
	ginC.JSON(http.StatusOK, BaseResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

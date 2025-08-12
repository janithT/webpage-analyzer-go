package responses

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var mu sync.Mutex

type BaseResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteError sends a standardized error response to the client
func WriteError(ginC *gin.Context, statusCode int, message string) {
	mu.Lock()
	defer mu.Unlock()
	ginC.JSON(statusCode, BaseResponse{
		Status:  "error",
		Message: message,
	})
	ginC.Abort()
}

// WriteSuccess sends a standardized success response with data
func WriteSuccess(ginC *gin.Context, message string, data interface{}) {
	mu.Lock()
	defer mu.Unlock()
	ginC.JSON(http.StatusOK, BaseResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

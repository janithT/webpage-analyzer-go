package responses

type ErrorResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponseWithStatus creates and returns a standardized error response
func ErrorResponseWithStatus(message string) ErrorResponse {
	return ErrorResponse{
		Status:  "error",
		Message: message,
		Data:    nil,
	}
}

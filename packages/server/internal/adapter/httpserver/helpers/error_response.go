package helpers

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func NewErrorResponse(status int, message string, err string) *ErrorResponse {
	return &ErrorResponse{
		Status:  status,
		Message: message,
		Error:   err,
	}
}

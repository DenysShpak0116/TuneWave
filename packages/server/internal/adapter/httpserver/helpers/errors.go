package helpers

import "fmt"

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(status int, format string, args ...interface{}) error {
	return &APIError{
		Status:  status,
		Message: fmt.Sprintf(format, args...),
	}
}

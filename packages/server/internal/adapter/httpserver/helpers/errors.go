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

func BadRequest(format string, args ...interface{}) error {
	return NewAPIError(400, format, args...)
}

func NotFound(format string, args ...interface{}) error {
	return NewAPIError(404, format, args...)
}

func InternalServerError(format string, args ...interface{}) error {
	return NewAPIError(500, format, args...)
}

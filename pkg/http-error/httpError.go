package http_error

import "fmt"

type HttpError struct {
	Code    int
	Message string
}

func (t *HttpError) Error() string {
	return fmt.Sprintf("code=%d, message: %s", t.Code, t.Message)
}

func NewHttpError(code int, message string) error {
	return &HttpError{
		Code:    code,
		Message: message,
	}
}

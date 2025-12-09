package model

import nethttp "net/http"

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{Code: code, Message: message}
}

var (
	ErrBadRequest          = NewHTTPError(nethttp.StatusBadRequest, "Bad Request")
	ErrUnauthorized        = NewHTTPError(nethttp.StatusUnauthorized, "Unauthorized")
	ErrNotFound            = NewHTTPError(nethttp.StatusNotFound, "Not Found")
	ErrConflict            = NewHTTPError(nethttp.StatusConflict, "Conflict")
	ErrInternalServerError = NewHTTPError(nethttp.StatusInternalServerError, "Internal Server Error")
)

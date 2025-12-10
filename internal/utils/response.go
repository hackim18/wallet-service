package utils

import (
	"errors"
	"net/http"
	"wallet-service/internal/model"

	"github.com/gin-gonic/gin"
)

func FailedResponse(ctx *gin.Context, statusCode int, message string, err error) model.WebResponse[any] {
	return model.WebResponse[any]{
		Errors: message,
	}
}

type HTTPError interface {
	error
	Status() int
	Message() string
	Unwrap() error
}

type appError struct {
	msg    string
	status int
	err    error
}

func (e *appError) Error() string   { return e.msg }
func (e *appError) Message() string { return e.msg }
func (e *appError) Status() int     { return e.status }
func (e *appError) Unwrap() error   { return e.err }

func Error(message string, status int, err error) error {
	return &appError{msg: message, status: status, err: err}
}

func HandleHTTPError(ctx *gin.Context, err error) {
	var httpErr HTTPError
	var modelHTTPError *model.HTTPError

	if errors.As(err, &modelHTTPError) {
		res := FailedResponse(ctx, modelHTTPError.Code, modelHTTPError.Message, nil)
		ctx.AbortWithStatusJSON(modelHTTPError.Code, res)
		return
	}

	if errors.As(err, &httpErr) {
		res := FailedResponse(ctx, httpErr.Status(), httpErr.Message(), httpErr.Unwrap())
		ctx.AbortWithStatusJSON(httpErr.Status(), res)
		return
	}

	res := FailedResponse(ctx, http.StatusInternalServerError, "Internal server error", err)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
}

func SuccessResponse[T any](ctx *gin.Context, statusCode int, message string, data T) model.WebResponse[T] {
	return model.WebResponse[T]{
		Data:    data,
		Message: message,
	}
}

func SuccessWithPaginationResponse[T any](
	ctx *gin.Context,
	statusCode int,
	message string,
	data []T,
	paging model.PageMetadata,
	documentationURL ...string,
) model.WebResponse[[]T] {
	return model.WebResponse[[]T]{
		Message: message,
		Data:    data,
		Paging:  &paging,
	}
}

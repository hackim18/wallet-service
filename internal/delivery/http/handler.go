package http

import (
	nethttp "net/http"

	"wallet-service/internal/model"

	"github.com/gin-gonic/gin"
)

func handleError(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	if httpErr, ok := err.(*model.HTTPError); ok {
		ctx.AbortWithStatusJSON(httpErr.Code, model.WebResponse[any]{Errors: httpErr.Message})
		return
	}

	ctx.AbortWithStatusJSON(nethttp.StatusInternalServerError, model.WebResponse[any]{Errors: err.Error()})
}

func Handler(h func(*gin.Context) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := h(ctx); err != nil {
			handleError(ctx, err)
		}
	}
}

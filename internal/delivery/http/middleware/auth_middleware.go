package middleware

import (
	nethttp "net/http"
	"wallet-service/internal/model"
	"wallet-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewAuth(userUserCase *usecase.UserUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &model.VerifyUserRequest{Token: ctx.GetHeader("Authorization")}
		userUserCase.Log.Debugf("Authorization : %s", request.Token)

		auth, err := userUserCase.Verify(ctx.Request.Context(), request)
		if err != nil {
			userUserCase.Log.Warnf("Failed find user by token : %+v", err)
			ctx.AbortWithStatusJSON(nethttp.StatusUnauthorized, model.WebResponse[any]{Errors: model.ErrUnauthorized.Message})
			return
		}

		userUserCase.Log.Debugf("User : %+v", auth.ID)
		ctx.Set("auth", auth)
		ctx.Next()
	}
}

func GetUser(ctx *gin.Context) *model.Auth {
	auth, _ := ctx.Get("auth")
	return auth.(*model.Auth)
}

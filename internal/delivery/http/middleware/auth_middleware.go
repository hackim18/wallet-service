package middleware

import (
	"net/http"
	"wallet-service/internal/constants"
	"wallet-service/internal/model"
	"wallet-service/internal/usecase"
	"wallet-service/internal/utils"

	"github.com/gin-gonic/gin"
)

func NewAuth(userUserCase *usecase.UserUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &model.VerifyUserRequest{Token: ctx.GetHeader("Authorization")}
		userUserCase.Log.Debugf("Authorization : %s", request.Token)

		auth, err := userUserCase.Verify(ctx.Request.Context(), request)
		if err != nil {
			res := utils.FailedResponse(ctx, http.StatusUnauthorized, constants.InvalidToken, err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		userUserCase.Log.Debugf("User : %+v", auth.UserID)
		ctx.Set("auth", auth)
		ctx.Next()
	}
}

func GetUser(ctx *gin.Context) *model.Auth {
	auth, _ := ctx.Get("auth")
	return auth.(*model.Auth)
}

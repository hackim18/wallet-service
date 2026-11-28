package route

import (
	"github.com/hackim18/wallet-service/internal/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *RouteConfig) RegisterCommonRoutes(app *gin.Engine) {
	welcomeHandler := func(ctx *gin.Context) {
		res := gin.H{"message": constants.WelcomeMessage}
		ctx.JSON(http.StatusOK, res)
	}

	app.GET("/", welcomeHandler)
	app.GET("/api", welcomeHandler)
	app.NoRoute(func(ctx *gin.Context) {
		res := gin.H{"message": constants.NotFound}
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
	})
}

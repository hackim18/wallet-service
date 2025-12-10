package route

import (
	"net/http"
	"wallet-service/internal/constants"

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

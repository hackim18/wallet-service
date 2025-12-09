package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *RouteConfig) RegisterCommonRoutes(app *gin.Engine) {
	welcomeHandler := func(ctx *gin.Context) {
		res := gin.H{"message": "Welcome to Wallet Service API"}
		ctx.JSON(http.StatusOK, res)
	}

	app.GET("/", welcomeHandler)
	app.GET("/api", welcomeHandler)
	app.NoRoute(func(ctx *gin.Context) {
		res := gin.H{"message": "API not found"}
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
	})
}

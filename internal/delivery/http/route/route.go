package route

import (
	"wallet-service/internal/delivery/http"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	Router         *gin.Engine
	UserController *http.UserController
	AuthMiddleware gin.HandlerFunc
}

func (c *RouteConfig) Setup() {
	api := c.Router.Group("/api")

	c.RegisterUserRoutes(api)
	c.RegisterCommonRoutes(c.Router)
}

package route

import "github.com/gin-gonic/gin"

func (c *RouteConfig) RegisterUserRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/users")

	user.POST("/register", c.UserController.Register)
	user.POST("/login", c.UserController.Login)
}

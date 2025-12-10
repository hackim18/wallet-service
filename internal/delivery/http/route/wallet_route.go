package route

import "github.com/gin-gonic/gin"

func (c *RouteConfig) RegisterWalletRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/wallets")

	user.GET("/balance", c.AuthMiddleware, c.WalletController.GetBalance)
}

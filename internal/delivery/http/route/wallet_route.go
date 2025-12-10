package route

import "github.com/gin-gonic/gin"

func (c *RouteConfig) RegisterWalletRoutes(rg *gin.RouterGroup) {
	wallet := rg.Group("/wallets")

	wallet.GET("", c.AuthMiddleware, c.WalletController.List)
	wallet.GET("/:walletId/balance", c.AuthMiddleware, c.WalletController.GetBalance)
	wallet.POST("/:walletId/withdraw", c.AuthMiddleware, c.WalletController.Withdraw)
}

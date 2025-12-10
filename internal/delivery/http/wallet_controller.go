package http

import (
	"net/http"
	"wallet-service/internal/constants"
	"wallet-service/internal/delivery/http/middleware"
	"wallet-service/internal/usecase"
	"wallet-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type WalletController struct {
	UseCase *usecase.WalletUseCase
	Log     *logrus.Logger
}

func NewWalletController(useCase *usecase.WalletUseCase, logger *logrus.Logger) *WalletController {
	return &WalletController{
		UseCase: useCase,
		Log:     logger,
	}
}

func (c *WalletController) GetBalance(ctx *gin.Context) {
	auth := middleware.GetUser(ctx)

	currency := ctx.Query("currency")

	response, err := c.UseCase.GetBalance(ctx.Request.Context(), auth.UserID, currency)
	if err != nil {
		c.Log.WithError(err).Warn("failed to get wallet balance")
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusOK, constants.WalletBalanceFetched, response)
	ctx.JSON(http.StatusOK, res)
}

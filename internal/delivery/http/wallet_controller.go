package http

import (
	"net/http"
	"strconv"
	"strings"
	"wallet-service/internal/constants"
	"wallet-service/internal/delivery/http/middleware"
	"wallet-service/internal/model"
	"wallet-service/internal/usecase"
	"wallet-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type WalletController struct {
	UseCase  *usecase.WalletUseCase
	Log      *logrus.Logger
	Validate *validator.Validate
}

func NewWalletController(useCase *usecase.WalletUseCase, logger *logrus.Logger, validate *validator.Validate) *WalletController {
	return &WalletController{
		UseCase:  useCase,
		Log:      logger,
		Validate: validate,
	}
}

func (c *WalletController) GetBalance(ctx *gin.Context) {
	auth := middleware.GetUser(ctx)

	walletIDStr := ctx.Param("walletId")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		utils.HandleHTTPError(ctx, utils.Error(constants.ErrInvalidIDFormat, http.StatusBadRequest, err))
		return
	}

	response, err := c.UseCase.GetBalance(ctx.Request.Context(), auth.UserID, walletID)
	if err != nil {
		c.Log.WithError(err).Warn("failed to get wallet balance")
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusOK, constants.WalletBalanceFetched, response)
	ctx.JSON(http.StatusOK, res)
}

func (c *WalletController) Withdraw(ctx *gin.Context) {
	auth := middleware.GetUser(ctx)

	walletIDStr := ctx.Param("walletId")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		utils.HandleHTTPError(ctx, utils.Error(constants.ErrInvalidIDFormat, http.StatusBadRequest, err))
		return
	}

	request := new(model.WalletWithdrawRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).Warn("failed to parse withdraw request")
		utils.HandleHTTPError(ctx, utils.Error(constants.FailedDataFromBody, http.StatusBadRequest, err))
		return
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warn("withdraw validation failed")
		msg := utils.TranslateValidationError(c.Validate, err)
		utils.HandleHTTPError(ctx, utils.Error(msg, http.StatusBadRequest, err))
		return
	}

	response, err := c.UseCase.Withdraw(ctx.Request.Context(), auth.UserID, walletID, request.Amount, request.Reference, request.Description)
	if err != nil {
		c.Log.WithError(err).Warn("failed to withdraw")
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusOK, constants.MsgWithdrawSuccess, response)
	ctx.JSON(http.StatusOK, res)
}

func (c *WalletController) List(ctx *gin.Context) {
	auth := middleware.GetUser(ctx)

	response, err := c.UseCase.List(ctx.Request.Context(), auth.UserID)
	if err != nil {
		c.Log.WithError(err).Warn("failed to list wallets")
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusOK, constants.WalletListFetched, response)
	ctx.JSON(http.StatusOK, res)
}

func (c *WalletController) Deposit(ctx *gin.Context) {
	auth := middleware.GetUser(ctx)

	walletIDStr := ctx.Param("walletId")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		utils.HandleHTTPError(ctx, utils.Error(constants.ErrInvalidIDFormat, http.StatusBadRequest, err))
		return
	}

	request := new(model.WalletDepositRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).Warn("failed to parse deposit request")
		utils.HandleHTTPError(ctx, utils.Error(constants.FailedDataFromBody, http.StatusBadRequest, err))
		return
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warn("deposit validation failed")
		msg := utils.TranslateValidationError(c.Validate, err)
		utils.HandleHTTPError(ctx, utils.Error(msg, http.StatusBadRequest, err))
		return
	}

	response, err := c.UseCase.Deposit(ctx.Request.Context(), auth.UserID, walletID, request.Amount, request.Reference, request.Description)
	if err != nil {
		c.Log.WithError(err).Warn("failed to deposit")
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusOK, constants.MsgDepositSuccess, response)
	ctx.JSON(http.StatusOK, res)
}

func (c *WalletController) ListTransactions(ctx *gin.Context) {
	auth := middleware.GetUser(ctx)

	walletIDStr := ctx.Param("walletId")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		utils.HandleHTTPError(ctx, utils.Error(constants.ErrInvalidIDFormat, http.StatusBadRequest, err))
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	request := model.WalletTransactionsListRequest{
		Page: page,
		Size: size,
		Type: strings.ToUpper(strings.TrimSpace(ctx.Query("type"))),
	}

	response, paging, err := c.UseCase.ListTransactions(ctx.Request.Context(), auth.UserID, walletID, request)
	if err != nil {
		c.Log.WithError(err).Warn("failed to list wallet transactions")
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessWithPaginationResponse(ctx, http.StatusOK, constants.WalletTransactionsFetched, response, paging)
	ctx.JSON(http.StatusOK, res)
}

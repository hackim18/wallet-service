package http

import (
	"net/http"
	"wallet-service/internal/constants"
	"wallet-service/internal/model"
	"wallet-service/internal/usecase"
	"wallet-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log      *logrus.Logger
	UseCase  *usecase.UserUseCase
	Validate *validator.Validate
}

func NewUserController(useCase *usecase.UserUseCase, logger *logrus.Logger, validate *validator.Validate) *UserController {
	return &UserController{
		Log:      logger,
		UseCase:  useCase,
		Validate: validate,
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	request := new(model.RegisterUserRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		utils.HandleHTTPError(ctx, utils.Error(constants.FailedDataFromBody, http.StatusBadRequest, err))
		return
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Validation failed : %+v", err)
		message := utils.TranslateValidationError(c.Validate, err)
		utils.HandleHTTPError(ctx, utils.Error(message, http.StatusBadRequest, err))
		return
	}

	response, err := c.UseCase.Create(ctx.Request.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusCreated, constants.UserRegistered, response)
	ctx.JSON(http.StatusCreated, res)
}

func (c *UserController) Login(ctx *gin.Context) {
	request := new(model.LoginUserRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		utils.HandleHTTPError(ctx, utils.Error(constants.FailedDataFromBody, http.StatusBadRequest, err))
		return
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Validation failed : %+v", err)
		message := utils.TranslateValidationError(c.Validate, err)
		utils.HandleHTTPError(ctx, utils.Error(message, http.StatusBadRequest, err))
		return
	}

	response, err := c.UseCase.Login(ctx.Request.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusOK, constants.UserLoggedIn, response)
	ctx.JSON(http.StatusOK, res)
}

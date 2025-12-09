package http

import (
	nethttp "net/http"
	"wallet-service/internal/model"
	"wallet-service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserController(useCase *usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	request := new(model.RegisterUserRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return
	}

	response, err := c.UseCase.Create(ctx.Request.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return
	}

	ctx.JSON(nethttp.StatusOK, model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Login(ctx *gin.Context) {
	request := new(model.LoginUserRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
	}

	response, err := c.UseCase.Login(ctx.Request.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
	}

	ctx.JSON(nethttp.StatusOK, model.WebResponse[*model.UserResponse]{Data: response})
}

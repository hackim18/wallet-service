package config

import (
	"wallet-service/internal/delivery/http"
	"wallet-service/internal/delivery/http/middleware"
	"wallet-service/internal/delivery/http/route"
	"wallet-service/internal/repository"
	"wallet-service/internal/usecase"
	"wallet-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	Router   *gin.Engine
	DB       *gorm.DB
	JWT      *utils.JWTHelper
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	walletRepository := repository.NewWalletRepository(config.Log)

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.JWT, userRepository, walletRepository)
	walletUseCase := usecase.NewWalletUseCase(config.DB, config.Log, walletRepository)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log, config.Validate)
	walletController := http.NewWalletController(walletUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)

	routeConfig := route.RouteConfig{
		Router:           config.Router,
		UserController:   userController,
		WalletController: walletController,
		AuthMiddleware:   authMiddleware,
	}
	routeConfig.Setup()
}
